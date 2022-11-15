package usbhid

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func sysfsReadAsString(dir string, entry string) (string, error) {
	b, err := os.ReadFile(filepath.Join(dir, entry))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func sysfsReadAsUint(dir string, entry string, base int, bitSize int) (uint64, error) {
	v, err := sysfsReadAsString(dir, entry)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(v, base, bitSize)
}

func sysfsReadAsHexUint16(dir string, entry string) (uint16, error) {
	v, err := sysfsReadAsUint(dir, entry, 16, 16)
	return uint16(v), err
}

func listDevices() ([]*Device, error) {
	rv := []*Device{}

	if err := filepath.Walk("/sys/bus/usb/devices", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == 0 || strings.Contains(info.Name(), ":") {
			return nil
		}

		d := &Device{}

		d.vendorId, err = sysfsReadAsHexUint16(path, "idVendor")
		if err != nil {
			return nil
		}

		d.productId, err = sysfsReadAsHexUint16(path, "idProduct")
		if err != nil {
			return nil
		}

		d.version, err = sysfsReadAsHexUint16(path, "bcdDevice")
		if err != nil {
			return nil
		}

		if m, err := sysfsReadAsString(path, "manufacturer"); err == nil {
			d.manufacturer = m
		}

		if p, err := sysfsReadAsString(path, "product"); err == nil {
			d.product = p
		}

		if s, err := sysfsReadAsString(path, "serial"); err == nil {
			d.serialNumber = s
		}

		f, err := filepath.Glob(filepath.Join(path, "*", "*", "hidraw", "hidraw[0-9]*"))
		if err != nil {
			return nil
		}
		if len(f) != 1 {
			return nil
		}

		d.path = filepath.Join("/dev", filepath.Base(f[0]))
		rv = append(rv, d)

		return nil
	}); err != nil {
		return nil, err
	}

	return rv, nil
}

func (d *Device) lock() error {
	if d.file == nil {
		return errors.New("usbhid: device is not open")
	}

	if err := syscall.Flock(int(d.file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err == syscall.EWOULDBLOCK {
		return ErrDeviceLocked
	} else {
		return err
	}
}
