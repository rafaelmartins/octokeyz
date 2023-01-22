// Copyright 2022-2023 Rafael G.Martins. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package b8 provides support for interacting with a b8 USB keypad.
package b8

const (
	// USB vendor identifier reported by b8
	USBVendorId uint16 = 0x16c0

	// USB product identifier reported by b8
	USBProductId uint16 = 0x05df

	// Major USB version number reported by b8 that is supported by this package
	// (this is the interface version, not a USB protocol or product version)
	USBVersion byte = 0x01

	// USB manufacturer name reported by b8
	USBManufacturer = "rgm.io"

	// USB product name reported by b8
	USBProduct = "b8"
)
