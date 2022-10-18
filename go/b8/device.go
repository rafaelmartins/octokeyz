package b8

import (
	"fmt"
	"os"
)

type Device struct {
	evdev    string
	file     *os.File
	handlers []*handler
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

func (d *Device) AddHandler(button ButtonType, buttonState ButtonStateType, fn HandlerFunc) {
	d.handlers = append(d.handlers, &handler{
		fn:          fn,
		button:      button,
		buttonState: buttonState,
	})
}

func (d *Device) Listen() error {
	if d.file == nil {
		return fmt.Errorf("b8: char device is not open")
	}

	for {
		events, err := newEvents(d.file)
		if err != nil {
			return err
		}

		for _, ev := range events {
			for _, h := range d.handlers {
				h.execute(ev)
			}
		}
	}
}
