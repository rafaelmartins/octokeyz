package usbhid

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrDeviceLocked    = errors.New("device is locked by another application")
	ErrDeviceIsOpen    = errors.New("device is open")
	ErrDeviceIsNotOpen = errors.New("device is not open")
)

type Device struct {
	path         string
	vendorId     uint16
	productId    uint16
	version      uint16
	manufacturer string
	product      string
	serialNumber string
	file         *os.File
	flock        *os.File
}

type DeviceFilterFunc func(*Device) bool

func ListDevices(f DeviceFilterFunc) ([]*Device, error) {
	devices, err := listDevices()
	if err != nil {
		return nil, err
	}

	if f == nil {
		return devices, nil
	}

	rv := []*Device{}
	for _, dev := range devices {
		if f(dev) {
			rv = append(rv, dev)
		}
	}
	return rv, nil
}

func (d *Device) Open() error {
	if d.file != nil {
		return fmt.Errorf("usbhid: %s: %w", d.path, ErrDeviceIsOpen)
	}

	f, err := os.OpenFile(d.path, os.O_RDWR, 0755)
	if err != nil {
		return nil
	}

	d.file = f

	return d.lock()
}

func (d *Device) IsOpen() bool {
	return d.file != nil
}

func (d *Device) Close() error {
	if d.file == nil {
		return nil
	}

	if err := d.file.Close(); err != nil {
		return err
	}
	d.file = nil

	if d.flock != nil {
		fn := d.flock.Name()
		if err := d.flock.Close(); err != nil {
			return err
		}
		d.flock = nil
		os.Remove(fn)
	}

	return nil
}

func (d *Device) Read(buf []byte) (int, error) {
	if d.file == nil {
		return 0, fmt.Errorf("usbhid: %s: %w", d.path, ErrDeviceIsNotOpen)
	}

	return d.file.Read(buf)
}

func (d *Device) Write(buf []byte) (int, error) {
	if d.file == nil {
		return 0, fmt.Errorf("usbhid: %s: %w", d.path, ErrDeviceIsNotOpen)
	}

	return d.file.Write(buf)
}

func (d *Device) Path() string {
	return d.path
}

func (d *Device) VendorId() uint16 {
	return d.vendorId
}

func (d *Device) ProductId() uint16 {
	return d.productId
}

func (d *Device) Version() uint16 {
	return d.version
}

func (d *Device) Manufacturer() string {
	return d.manufacturer
}

func (d *Device) Product() string {
	return d.product
}

func (d *Device) SerialNumber() string {
	return d.serialNumber
}