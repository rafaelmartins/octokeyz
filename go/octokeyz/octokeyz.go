// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: GPL-2.0

// Package octokeyz provides support for interacting with an octokeyz USB macropad.
package octokeyz

const (
	// USB vendor identifier reported by octokeyz
	USBVendorId uint16 = 0x1d50

	// USB product identifier reported by octokeyz
	USBProductId uint16 = 0x6184

	// Major USB version number reported by octokeyz that is supported by this package
	// (this is the interface version, not a USB protocol or product version)
	USBVersion byte = 0x01

	// USB manufacturer name reported by octokeyz
	USBManufacturer = "rgm.io"

	// USB product name reported by octokeyz
	USBProduct = "octokeyz"
)
