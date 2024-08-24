// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: GPL-2.0

package octokeyz

import (
	"errors"
	"fmt"
	"time"

	"github.com/rafaelmartins/usbhid"
)

// Errors returned from octokeyz package may be tested against these errors
// with errors.Is.
var (
	ErrButtonInvalid                     = errors.New("button is not valid")
	ErrButtonHandlerInvalid              = errors.New("button handler is not valid")
	ErrDeviceDisplayNumberOfLinesInvalid = errors.New("device firmware reported an incompatible number of display lines")
	ErrDeviceDisplayNotSupported         = errors.New("device hardware does not includes a display")
	ErrDeviceEnumerationFailed           = usbhid.ErrDeviceEnumerationFailed
	ErrDeviceFailedToClose               = usbhid.ErrDeviceFailedToClose
	ErrDeviceFailedToOpen                = usbhid.ErrDeviceFailedToOpen
	ErrDeviceFirmwareVersionIncompatible = errors.New("device firmware version is not compatible")
	ErrDeviceIsClosed                    = usbhid.ErrDeviceIsClosed
	ErrDeviceIsOpen                      = usbhid.ErrDeviceIsOpen
	ErrDeviceLocked                      = usbhid.ErrDeviceLocked
	ErrDeviceMoreThanOne                 = errors.New("more than one device found")
	ErrDeviceNotFound                    = errors.New("device not found")
	ErrGetFeatureReportFailed            = usbhid.ErrGetFeatureReportFailed
	ErrGetInputReportFailed              = usbhid.ErrGetInputReportFailed
	ErrReportBufferOverflow              = usbhid.ErrReportBufferOverflow
	ErrSetFeatureReportFailed            = usbhid.ErrSetFeatureReportFailed
	ErrSetOutputReportFailed             = usbhid.ErrSetOutputReportFailed
)

// Device is an opaque structure that represents an octokeyz USB macropad device
// connected to the computer.
type Device struct {
	dev                 *usbhid.Device
	listen              chan bool
	buttons             map[ButtonID]*Button
	data                byte
	legacyLedState      bool
	withDisplay         bool
	displayCharsPerLine byte
}

// LedState represents a state to set the octokeyz USB macropad led to.
type LedState byte

const (
	// LedOn sets the led on.
	LedOn LedState = iota + 1

	// LedFlash sets the led to flash on for a short time and go off.
	LedFlash

	// LedSlowBlink sets the led to blink slowly.
	LedSlowBlink

	// LedFastBlink sets the led to blink fastly.
	LedFastBlink

	// LedOff sets the led off.
	LedOff
)

func wrapErr(err error) error {
	if err != nil {
		return fmt.Errorf("octokeyz: %w", err)
	}
	return nil
}

// Enumerate lists the octokeyz USB macropads connected to the computer.
func Enumerate() ([]*Device, error) {
	devices, err := usbhid.Enumerate(func(d *usbhid.Device) bool {
		return d.VendorId() == USBVendorId && d.ProductId() == USBProductId
	})
	if err != nil {
		return nil, wrapErr(err)
	}

	rv := []*Device{}
	for _, dev := range devices {
		rv = append(rv, &Device{
			dev: dev,
		})
	}
	return rv, nil
}

// GetDevice returns an octokeyz USB macropad found connected to the machine that matches the
// provided serial number. If serial number is empty and only one device is connected,
// this device is returned, otherwise an error is returned.
func GetDevice(serialNumber string) (*Device, error) {
	devices, err := Enumerate()
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		if serialNumber != "" {
			return nil, fmt.Errorf("octokeyz: %w [%q]", ErrDeviceNotFound, serialNumber)
		}
		return nil, wrapErr(ErrDeviceNotFound)
	}

	if serialNumber == "" {
		if len(devices) == 1 {
			return devices[0], nil
		}

		sn := []string{}
		for _, dev := range devices {
			sn = append(sn, dev.SerialNumber())
		}

		return nil, fmt.Errorf("octokeyz: %w %q", ErrDeviceMoreThanOne, sn)
	}

	for _, dev := range devices {
		if dev.SerialNumber() == serialNumber {
			return dev, nil
		}
	}

	return nil, fmt.Errorf("octokeyz: %w [%q]", ErrDeviceNotFound, serialNumber)
}

// Open opens the octokeyz USB macropad for usage.
func (d *Device) Open() error {
	if d.dev == nil {
		return wrapErr(ErrDeviceNotFound)
	}

	if byte(d.dev.Version()>>8) != USBVersion {
		return fmt.Errorf("octokeyz: %w (expected %02Xxx, got %04X)", ErrDeviceFirmwareVersionIncompatible, USBVersion, d.dev.Version())
	}

	if byte(d.dev.Version()) < 1 {
		d.legacyLedState = true
	}

	if err := d.dev.Open(true); err != nil {
		return wrapErr(err)
	}
	d.listen = make(chan bool)

	if buf, err := d.dev.GetFeatureReport(1); err == nil {
		d.withDisplay = buf[0]&(1<<0) != 0

		if d.withDisplay {
			buf, err = d.dev.GetFeatureReport(2)
			if err != nil {
				return wrapErr(err)
			}

			if buf[0] != 8 {
				return fmt.Errorf("octokeyz: %w [%d]", ErrDeviceDisplayNumberOfLinesInvalid, buf[0])
			}

			d.displayCharsPerLine = buf[1]

			if err := d.DisplayClear(); err != nil {
				return wrapErr(err)
			}
		}
	}
	return nil
}

// Close closes the octokeyz USB macropad.
func (d *Device) Close() error {
	if d.dev == nil {
		return wrapErr(ErrDeviceNotFound)
	}

	d.DisplayClear()
	d.Led(LedOff)
	close(d.listen)
	return wrapErr(d.dev.Close())
}

// AddHandler registers a ButtonHandler callback to be called whenever the given
// button is pressed.
func (d *Device) AddHandler(button ButtonID, fn ButtonHandler) error {
	if fn == nil {
		return wrapErr(ErrButtonHandlerInvalid)
	}

	if d.buttons == nil {
		d.buttons = newButtons()
	}

	if btn, ok := d.buttons[button]; ok {
		btn.addHandler(fn)
		return nil
	}
	return wrapErr(ErrButtonInvalid)
}

// Listen listens to button press events from the macropad and calls ButtonHandler
// callbacks as required.
//
// errCh is an error channel to receive errors from the button handlers. If set
// to a nil channel, errors are sent to standard logger. Errors are sent
// non-blocking.
func (d *Device) Listen(errCh chan error) error {
	if d.dev == nil {
		return wrapErr(ErrDeviceNotFound)
	}

	for {
		select {
		case <-d.listen:
			return nil
		default:
			if d.listen == nil {
				return nil
			}
		}

		id, buf, err := d.dev.GetInputReport()
		if err != nil {
			return wrapErr(err)
		}

		if id != 1 {
			continue
		}

		t := time.Now()

		if buf[0] == d.data {
			continue
		}

		for j := 0; j < 8; j++ {
			if v := buf[0] & (1 << j); v != (d.data & (1 << j)) {
				if btn, ok := d.buttons[BUTTON_1+ButtonID(j)]; ok {
					if v > 0 {
						btn.press(t, errCh)
					} else {
						btn.release(t)
					}
				}
			}
		}

		d.data = buf[0]
	}
}

// Led sets the state of the octokeyz USB macropad led.
func (d *Device) Led(state LedState) error {
	if d.dev == nil {
		return wrapErr(ErrDeviceNotFound)
	}

	if d.legacyLedState {
		state--
	}

	return wrapErr(d.dev.SetOutputReport(1, []byte{byte(state)}))
}

// SerialNumber returns the serial number of the octokeyz USB macropad.
func (d *Device) SerialNumber() string {
	return d.dev.SerialNumber()
}
