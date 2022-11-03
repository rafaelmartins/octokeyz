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

func (d *Device) Open() error {
	if d.dev == nil {
		return errors.New("b8: device not defined")
	}

	return d.dev.Open()
}

func (d *Device) Close() error {
	if d.dev == nil {
		return nil
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
	if d.dev == nil || !d.dev.IsOpen() {
		return errors.New("b8: device is not open")
	}

	buf := make([]byte, inputReportLen*64)

	for {
		n, err := d.dev.Read(buf)
		if err != nil {
			return err
		}

		if n%inputReportLen != 0 {
			return errors.New("b8: failed to read hid report")
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
