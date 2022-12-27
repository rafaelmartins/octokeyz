package b8

import (
	"errors"
	"fmt"
	"time"

	"github.com/rafaelmartins/b8/go/b8/internal/usbhid"
)

const (
	inputReportLen  = 2
	outputReportLen = 2
	reportID        = 1
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

func ListDevices() ([]*Device, error) {
	devices, err := usbhid.ListDevices(func(d *usbhid.Device) bool {
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
		if (dev.Version() >> 8) != (USBVersion >> 8) {
			return nil, fmt.Errorf("b8: device version is not compatible, please upgrade: %s: 0x%04x", dev.Path(), dev.Version())
		}

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
		return nil, fmt.Errorf("b8: %q: %w", serialNumber, ErrDeviceNotFound)
	}

	if serialNumber == "" {
		if len(devices) == 1 {
			return devices[0], nil
		}

		sn := []string{}
		for _, dev := range devices {
			sn = append(sn, dev.SerialNumber())
		}
		return nil, fmt.Errorf("b8: %w: %q", ErrDeviceMoreThanOne, sn)
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
		return ErrDeviceNotFound
	}

	return d.dev.Open()
}

func (d *Device) Close() error {
	if d.dev == nil {
		return ErrDeviceNotFound
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
		return ErrDeviceNotFound
	}

	buf := make([]byte, inputReportLen*64)

	for {
		n, err := d.dev.Read(buf)
		if err != nil {
			return fmt.Errorf("b8: %w: %s", ErrDeviceReadFailed, err)
		}

		if n%inputReportLen != 0 {
			return fmt.Errorf("b8: %w: bad read size", ErrDeviceReadFailed)
		}

		t := time.Now()

		for i := 0; i < int(n); i += inputReportLen {
			if buf[i] != reportID || buf[i+1] == d.data {
				continue
			}

			for j := 0; j < 8; j++ {
				if v := buf[i+1] & (1 << j); v != (d.data & (1 << j)) {
					if btn, ok := d.buttons[ButtonID(j)]; ok {
						if v > 0 {
							btn.press(t)
						} else {
							btn.release(t)
						}
					}
				}
			}

			d.data = buf[i+1]
		}
	}
}

func (d *Device) Led(state LedState) error {
	return led(d, state)
}

func (d *Device) SerialNumber() string {
	return d.dev.SerialNumber()
}
