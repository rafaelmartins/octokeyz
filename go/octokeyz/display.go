// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: GPL-2.0

package octokeyz

import (
	"encoding/binary"
	"time"
)

// DisplayLine represents the identifier of a display line number.
type DisplayLine byte

// An octokeyz USB macropad may contain up to 8 lines
const (
	DisplayLine1 DisplayLine = iota + 1
	DisplayLine2
	DisplayLine3
	DisplayLine4
	DisplayLine5
	DisplayLine6
	DisplayLine7
	DisplayLine8
)

// DisplayLine represents the identifier of a display line horizontal alignment.
type DisplayLineAlign byte

// A display line may have its content horizontally aligned left, right or center.
const (
	DisplayLineAlignLeft DisplayLineAlign = iota + 1
	DisplayLineAlignRight
	DisplayLineAlignCenter
)

// DisplayLine draws the provided string to the provided display line with the provided
// horizontal alignment. An error may be returned.
func (d *Device) DisplayLine(line DisplayLine, str string, align DisplayLineAlign) error {
	if !d.withDisplay {
		return wrapErr(ErrDeviceDisplayNotSupported)
	}

	end := len(str)
	if cpl := int(d.displayCharsPerLine); end > cpl {
		end = cpl
	}
	data := append([]byte(str[:end]), make([]byte, int(d.displayCharsPerLine)-end)...)
	data = append([]byte{byte(line) - 1, byte(align)}, data...)

	return wrapErr(d.dev.SetOutputReport(2, data))
}

// DisplayClearLine erases the provided display line. An error may be returned.
func (d *Device) DisplayClearLine(line DisplayLine) error {
	return d.DisplayLine(line, "", DisplayLineAlignLeft)
}

// DisplayClearWithDelay erases the whole display after the provided delay (milliseconds
// resolution). An error may be returned if the display does not support this feature.
// If a new line is drawn while the delay is ongoing the display is cleared immediately.
func (d *Device) DisplayClearWithDelay(delay time.Duration) error {
	if !d.withDisplay {
		return wrapErr(ErrDeviceDisplayNotSupported)
	}
	if !d.withDisplayClearWithDelay {
		return wrapErr(ErrDeviceDisplayClearWithDelayNotSupported)
	}

	return wrapErr(d.dev.SetOutputReport(3, binary.LittleEndian.AppendUint16(nil, uint16(delay/time.Millisecond))))
}

// DisplayClear erases the whole display. An error may be returned.
func (d *Device) DisplayClear() error {
	if d.withDisplayClearWithDelay {
		return d.DisplayClearWithDelay(0)
	}

	for i := DisplayLine1; i <= DisplayLine8; i++ {
		if err := d.DisplayClearLine(i); err != nil {
			return err
		}
	}
	return nil
}

// GetDisplayCharsPerLine returns how many characters can fit in a display line
// without overflowing.
func (d *Device) GetDisplayCharsPerLine() byte {
	return d.displayCharsPerLine
}
