package b8

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	idVendor     = 0x16c0
	idProduct    = 0x05df
	manufacturer = "rgm.io"
	product      = "b8"
)

func ListDevices() ([]*Device, error) {
	devices := []*Device{}

	if err := filepath.Walk("/sys/bus/usb/devices", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == 0 || strings.Contains(info.Name(), ":") {
			return nil
		}

		idv, err := sysfsReadAsHexUint16(path, "idVendor")
		if err != nil {
			return err
		}
		if idv != idVendor {
			return nil
		}

		idp, err := sysfsReadAsHexUint16(path, "idProduct")
		if err != nil {
			return err
		}
		if idp != idProduct {
			return nil
		}

		m, err := sysfsReadAsString(path, "manufacturer")
		if err != nil {
			return err
		}
		if m != manufacturer {
			return nil
		}

		p, err := sysfsReadAsString(path, "product")
		if err != nil {
			return err
		}
		if p != product {
			return nil
		}

		ev, err := filepath.Glob(filepath.Join(path, "*", "*", "input", "input[0-9]*", "event[0-9]*"))
		if err != nil {
			return err
		}
		if len(ev) != 1 {
			return fmt.Errorf("b8: more than one evdev char device available")
		}

		devices = append(devices, &Device{
			evdev: filepath.Join("/dev", "input", filepath.Base(ev[0])),
		})

		return nil
	}); err != nil {
		return nil, err
	}

	return devices, nil
}
