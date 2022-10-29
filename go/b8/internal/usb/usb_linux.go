package usb

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	kernelEventSize = int(unsafe.Sizeof(kernelEvent{}))
)

type kernelEvent struct {
	time_ struct {
		sec  int64
		usec int64
	}
	type_ uint16
	code  uint16
	value int32
}

type platformContext struct {
	file *os.File
}

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

		if m, err := sysfsReadAsString(path, "manufacturer"); err == nil {
			d.manufacturer = m
		}

		if p, err := sysfsReadAsString(path, "product"); err == nil {
			d.product = p
		}

		if s, err := sysfsReadAsString(path, "serial"); err == nil {
			d.serialNumber = s
		}

		f, err := filepath.Glob(filepath.Join(path, "*", "*", "input", "input[0-9]*", "event[0-9]*"))
		if err != nil {
			return nil
		}
		if len(f) != 1 {
			return nil
		}

		d.path = filepath.Join("/dev", "input", filepath.Base(f[0]))
		rv = append(rv, d)

		return nil
	}); err != nil {
		return nil, err
	}

	return rv, nil
}

func (d *Device) Open() error {
	if d.open {
		return errors.New("usb: device is open")
	}

	f, err := os.OpenFile(d.path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	d.pctx.file = f
	d.open = true

	return nil
}

func (d *Device) Close() error {
	if !d.open {
		return nil
	}

	if err := d.pctx.file.Close(); err != nil {
		return err
	}

	d.pctx.file = nil
	d.open = false

	return nil
}

func (d *Device) Read() ([]*Event, error) {
	if !d.open {
		return nil, errors.New("usb: device is not open")
	}

	buf := make([]byte, 64*kernelEventSize)
	n, err := d.pctx.file.Read(buf)
	if err != nil {
		return nil, err
	}
	if n%kernelEventSize != 0 {
		return nil, errors.New("usb: failed to read hid report")
	}

	rv := []*Event{}

	for i := 0; i < n/kernelEventSize; i++ {
		ev := (*kernelEvent)(unsafe.Pointer(&buf[i*kernelEventSize]))
		rv = append(rv, &Event{
			Time:  time.Unix(ev.time_.sec, ev.time_.usec*1000),
			Type:  ev.type_,
			Code:  ev.code,
			Value: ev.value,
		})
	}

	return rv, nil
}

func (d *Device) Write(e *Event) error {
	if !d.open {
		return errors.New("usb: device is not open")
	}

	buf := make([]byte, kernelEventSize)
	ev := (*kernelEvent)(unsafe.Pointer(&buf[0]))
	ev.time_.sec = e.Time.Unix()
	ev.time_.usec = e.Time.UnixMicro() % 1000000
	ev.type_ = e.Type
	ev.code = e.Code
	ev.value = e.Value

	n, err := d.pctx.file.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return errors.New("usb: failed to write event")
	}

	return nil
}
