package b8

import (
	"fmt"
)

type LedState byte

const (
	LedOn = iota
	LedFlash
	LedSlowBlink
	LedFastBlink
	LedOff
)

func led(d *Device, s LedState) error {
	if d.dev == nil {
		return ErrDeviceNotFound
	}

	n, err := d.dev.Write([]byte{reportID, byte(s)})
	if err != nil {
		return fmt.Errorf("b8: %w: %s", ErrDeviceWriteFailed, err)
	}
	if n != outputReportLen {
		return fmt.Errorf("b8: %w: bad write size", ErrDeviceWriteFailed)
	}

	return nil
}
