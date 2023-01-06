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

var (
	ErrDeviceNotFound    = errors.New("device not found")
	ErrDeviceMoreThanOne = errors.New("more than one device found")
	ErrDeviceReadFailed  = errors.New("failed to read hid report")
	ErrDeviceWriteFailed = errors.New("failed to write hid report")
)

type Device struct {
	dev     *usbhid.Device
	buttons map[ButtonID]*Button
	data    byte
}

type LedState byte

const (
	LedOn = iota
	LedFlash
	LedSlowBlink
	LedFastBlink
	LedOff
)

func ListDevices() ([]*Device, error) {
	devices, err := usbhid.Enumerate(func(d *usbhid.Device) bool {
		if d.VendorId() != USBVendorId {
			return false
		}

		if d.ProductId() != USBProductId {
			return false
		}

		if d.Manufacturer() != USBManufacturer {
			return false
		}

		if d.Product() != USBProduct {
			return false
		}

		return true
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

func GetDevice(serialNumber string) (*Device, error) {
	devices, err := ListDevices()
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

func (d *Device) Open() error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	if (d.dev.Version() >> 8) != (USBVersion >> 8) {
		return fmt.Errorf("b8: device version is not compatible, please upgrade: %s: 0x%04x", d.dev.Path(), d.dev.Version())
	}

	return d.dev.Open(true)
}

func (d *Device) Close() error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	return d.dev.Close()
}

func (d *Device) AddHandler(button ButtonID, fn ButtonHandler) {
	if fn == nil {
		return
	}

	if d.buttons == nil {
		d.buttons = newButtons()
	}

	if btn, ok := d.buttons[button]; ok {
		btn.addHandler(fn)
	}
}

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
				if btn, ok := d.buttons[ButtonID(j)]; ok {
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

func (d *Device) Led(state LedState) error {
	if d.dev == nil {
		return fmt.Errorf("b8: %w", ErrDeviceNotFound)
	}

	if err := d.dev.SetOutputReport(reportID, []byte{byte(state)}); err != nil {
		return fmt.Errorf("b8: %w: %s", ErrDeviceWriteFailed, err)
	}
	return nil
}

func (d *Device) SerialNumber() string {
	return d.dev.SerialNumber()
}
