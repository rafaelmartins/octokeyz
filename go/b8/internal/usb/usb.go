package usb

import "time"

type Device struct {
	path         string
	open         bool
	vendorId     uint16
	productId    uint16
	manufacturer string
	product      string
	serialNumber string
	pctx         platformContext
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

func (d *Device) IsOpen() bool {
	return d.open
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

func (d *Device) Manufacturer() string {
	return d.manufacturer
}

func (d *Device) Product() string {
	return d.product
}

func (d *Device) SerialNumber() string {
	return d.serialNumber
}

type Event struct {
	Time  time.Time
	Type  uint16
	Code  uint16
	Value int32
}
