package b8

import (
	"errors"
	"time"

	"github.com/rafaelmartins/b8/go/b8/internal/usb"
)

const (
	evKey    = 1
	evLed    = 17
	btnMacro = 0x0290
	ledMisc  = 0x08
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

	return d.dev.Open()
}

func (d *Device) Close() error {
	if d.dev == nil {
		return nil
	}

	return d.dev.Close()
}

func (d *Device) AddHandler(button ButtonID, fn ButtonHandler) {
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

	for {
		events, err := d.dev.Read()
		if err != nil {
			return err
		}

		for _, ev := range events {
			if ev.Type != evKey {
				continue
			}

			if btn, ok := d.buttons[ButtonID(ev.Code-btnMacro)]; ok {
				if ev.Value > 0 {
					btn.press(ev.Time)
				} else {
					btn.release(ev.Time)
				}
			}
		}
	}
}

func (d *Device) led(v int32) error {
	if d.dev == nil || !d.dev.IsOpen() {
		return errors.New("b8: device is not open")
	}

	return d.dev.Write(&usb.Event{
		Time:  time.Now(),
		Type:  evLed,
		Code:  ledMisc,
		Value: v,
	})
}

func (d *Device) LedOn() error {
	return d.led(1)
}

func (d *Device) LedOff() error {
	return d.led(0)
}
