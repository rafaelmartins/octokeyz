package b8

import (
	"errors"
	"os"
)

type Device struct {
	evdev   string
	file    *os.File
	buttons map[ButtonID]*Button
}

func (d *Device) Open() error {
	if d.file != nil {
		return nil
	}

	f, err := os.Open(d.evdev)
	if err != nil {
		return err
	}
	d.file = f

	d.buttons = newButtons()

	return nil
}

func (d *Device) Close() error {
	if d.file == nil {
		return nil
	}

	err := d.file.Close()
	d.file = nil
	return err
}

func (d *Device) AddHandler(button ButtonID, fn ButtonHandler) {
	d.buttons[button].addHandler(fn)
}

func (d *Device) Listen() error {
	if d.file == nil {
		return errors.New("b8: char device is not open")
	}

	for {
		events, err := newEvents(d.file)
		if err != nil {
			return err
		}

		for _, ev := range events {
			if ev.pressed {
				d.buttons[ev.button].press(ev.etime)
			} else {
				d.buttons[ev.button].release(ev.etime)
			}
		}
	}
}
