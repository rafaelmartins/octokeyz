// Copyright 2022-2023 Rafael G.Martins. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package b8

import (
	"errors"
	"fmt"
	"time"

	"github.com/rafaelmartins/usbhid"
)

const (
	reportID = 1
)

// Errors returned from b8 package may be tested against these errors
// with errors.Is.
var (
	ErrButtonInvalid        = errors.New("button is not valid")
	ErrButtonHandlerInvalid = errors.New("button handler is not valid")
	ErrDeviceNotFound       = errors.New("device not found")
	ErrDeviceMoreThanOne    = errors.New("more than one device found")
	ErrDeviceReadFailed     = errors.New("failed to read hid report")
	ErrDeviceWriteFailed    = errors.New("failed to write hid report")
)

// Device is an opaque structure that represents a b8 USB keypad device
// connected to the computer.
type Device struct {
	dev     *usbhid.Device
	buttons map[ButtonID]*Button
	data    byte
}

// LedState represents a state to set the b8 USB keypad led to.
type LedState byte

const (
	// LedOn sets the led on.
	LedOn = iota

	// LedFlash sets the led to flash on for a short time and go off.
	LedFlash

	// LedSlowBlink sets the led to blink slowly.
	LedSlowBlink

	// LedFastBlink sets the led to blink fastly.
	LedFastBlink

	// LedOff sets the led off.
	LedOff
)

// Enumerate lists the b8 USB keypads connected to the computer.
func Enumerate() ([]*Device, error) {
	devices, err := usbhid.Enumerate(func(d *usbhid.Device) bool {
		switch d.VendorId() {
		case USBVendorId:
			return d.ProductId() == USBProductId

		case 0x16c0: // old v-usb shared vid
			if d.ProductId() != 0x05df { // old v-usb shared hid pid
				return false
			}
			if d.Manufacturer() != USBManufacturer {
				return false
			}
			if d.Product() != USBProduct {
				return false
			}
			return true

		default:
			return false
		}
	})
	if err != nil {
		return nil, err
	}

	rv := []*Device{}
	for _, dev := range devices {
		rv = append(rv, &Device{
			dev: dev,
		})
	}
	return rv, nil
}

// GetDevice returns a b8 USB keypad found connected to the machine that matches the
// provided serial number. If serial number is empty and only one device is connected,
// this device is returned, otherwise an error is returned.
func GetDevice(serialNumber string) (*Device, error) {
	devices, err := Enumerate()
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		if serialNumber != "" {
			return nil, fmt.Errorf("b8: %q: %w", serialNumber, ErrDeviceNotFound)
		}
		return nil, fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	if serialNumber == "" {
		if len(devices) == 1 {
			return devices[0], nil
		}

		sn := []string{}
		for _, dev := range devices {
			sn = append(sn, dev.SerialNumber())
		}

		return nil, fmt.Errorf("b8: %q: %w", sn, ErrDeviceMoreThanOne)
	}

	for _, dev := range devices {
		if dev.SerialNumber() == serialNumber {
			return dev, nil
		}
	}

	return nil, fmt.Errorf("b8: %q: %w", serialNumber, ErrDeviceNotFound)
}

// Open opens the b8 USB keypad for usage.
func (d *Device) Open() error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	if byte(d.dev.Version()>>8) != USBVersion {
		return fmt.Errorf("b8: device version is not compatible, please upgrade: %s: 0x%04x", d.dev.Path(), d.dev.Version())
	}

	return d.dev.Open(true)
}

// Close closes the b8 USB keypad.
func (d *Device) Close() error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	return d.dev.Close()
}

// AddHandler registers a ButtonHandler callback to be called whenever the given
// button is pressed.
func (d *Device) AddHandler(button ButtonID, fn ButtonHandler) error {
	if fn == nil {
		return ErrButtonHandlerInvalid
	}

	if d.buttons == nil {
		d.buttons = newButtons()
	}

	if btn, ok := d.buttons[button]; ok {
		btn.addHandler(fn)
		return nil
	}
	return ErrButtonInvalid
}

// Listen listens to button press events from the keypad and calls ButtonHandler
// callbacks as required.
func (d *Device) Listen() error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	for {
		id, buf, err := d.dev.GetInputReport()
		if err != nil {
			return fmt.Errorf("b8: %w: %s", ErrDeviceReadFailed, err)
		}

		if id != reportID {
			return fmt.Errorf("b8: %w: bad input report id: %d", ErrDeviceReadFailed, id)
		}

		t := time.Now()

		if buf[0] == d.data {
			continue
		}

		for j := 0; j < 8; j++ {
			if v := buf[0] & (1 << j); v != (d.data & (1 << j)) {
				if btn, ok := d.buttons[BUTTON_1+ButtonID(j)]; ok {
					if v > 0 {
						btn.press(t)
					} else {
						btn.release(t)
					}
				}
			}
		}

		d.data = buf[0]
	}
}

// Led sets the state of the b8 USB keypad led.
func (d *Device) Led(state LedState) error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	if err := d.dev.SetOutputReport(reportID, []byte{byte(state)}); err != nil {
		return fmt.Errorf("b8: %w: %s", ErrDeviceWriteFailed, err)
	}
	return nil
}

// SerialNumber returns the serial number of the b8 USB keypad.
func (d *Device) SerialNumber() string {
	return d.dev.SerialNumber()
}
