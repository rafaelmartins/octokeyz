package b8

import (
	"errors"

	"github.com/rafaelmartins/b8/go/b8/internal/usb"
)

const (
	btnMacro = 0x290
)

type Device struct {
	dev     *usb.Device
	buttons map[ButtonID]*Button
}

func ListDevices() ([]*Device, error) {
	devices, err := usb.ListDevices(func(d *usb.Device) bool {
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

func (d *Device) Open() error {
	if d.dev == nil {
		return errors.New("b8: device not defined")
	}

	if err := d.dev.Open(); err != nil {
		return err
	}

	d.buttons = newButtons()
	return nil
}

func (d *Device) Close() error {
	if d.dev == nil {
		return nil
	}

	return d.dev.Close()
}

func (d *Device) AddHandler(button ButtonID, fn ButtonHandler) {
	if d.buttons != nil {
		if btn, ok := d.buttons[button]; ok {
			btn.addHandler(fn)
		}
	}
}

func (d *Device) Listen() error {
	if d.dev == nil || !d.dev.IsOpen() {
		return errors.New("b8: device is not open")
	}

	for {
		events, err := d.dev.Read()
		if err != nil {
			return err
		}

		for _, ev := range events {
			if ev.IsPressed() {
				d.buttons[ButtonID(ev.Key()-btnMacro)].press(ev.Time())
			} else {
				d.buttons[ButtonID(ev.Key()-btnMacro)].release(ev.Time())
			}
		}
	}
}
