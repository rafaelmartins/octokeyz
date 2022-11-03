package b8

import "errors"

type LedState byte

const (
	LedOn = iota
	LedFlash
	LedSlowBlink
	LedFastBlink
	LedOff
)

func led(d *Device, s LedState) error {
	if d.dev == nil || !d.dev.IsOpen() {
		return errors.New("b8: device is not open")
	}

	n, err := d.dev.Write([]byte{reportID, byte(s)})
	if err != nil {
		return err
	}
	if n != outputReportLen {
		return errors.New("b8: failed to write hid report")
	}

	return nil
}
