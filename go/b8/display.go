// Copyright 2022-2024 Rafael G. Martins. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package b8

import (
	"fmt"
)

// DisplayLine represents the identifier of a display line number.
type DisplayLine byte

// A b8 USB keypad may contain up to 8 lines
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
		return fmt.Errorf("b8: %w", ErrDisplayNotSupported)
	}

	end := len(str)
	if cpl := int(d.displayCharsPerLine); end > cpl {
		end = cpl
	}
	data := append([]byte(str[:end]), make([]byte, int(d.displayCharsPerLine)-end)...)
	data = append([]byte{byte(line) - 1, byte(align)}, data...)

	return d.dev.SetOutputReport(2, data)
}

// DisplayClearLine erases the provided display line. An error may be returned.
func (d *Device) DisplayClearLine(line DisplayLine) error {
	return d.DisplayLine(line, "", DisplayLineAlignLeft)
}

// DisplayClear erases the whole display. An error may be returned.
func (d *Device) DisplayClear() error {
	for i := DisplayLine1; i <= DisplayLine8; i++ {
		if err := d.DisplayClearLine(i); err != nil {
			return err
		}
	}
	return nil
}
