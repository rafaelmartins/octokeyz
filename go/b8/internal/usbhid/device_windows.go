package usbhid

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	sDIGCF_PRESENT         = 0x02
	sDIGCF_DEVICEINTERFACE = 0x10
)

var (
	setupapi                         = syscall.NewLazyDLL("setupapi.dll")
	setupDiDestroyDeviceInfoList     = setupapi.NewProc("SetupDiDestroyDeviceInfoList")
	setupDiEnumDeviceInterfaces      = setupapi.NewProc("SetupDiEnumDeviceInterfaces")
	setupDiGetClassDevsA             = setupapi.NewProc("SetupDiGetClassDevsA")
	setupDiGetDeviceInterfaceDetailA = setupapi.NewProc("SetupDiGetDeviceInterfaceDetailA")
)

var (
	hid                        = syscall.NewLazyDLL("hid.dll")
	hidD_GetAttributes         = hid.NewProc("HidD_GetAttributes")
	hidD_GetHidGuid            = hid.NewProc("HidD_GetHidGuid")
	hidD_GetManufacturerString = hid.NewProc("HidD_GetManufacturerString")
	hidD_GetProductString      = hid.NewProc("HidD_GetProductString")
	hidD_GetSerialNumberString = hid.NewProc("HidD_GetSerialNumberString")
)

type gGUID struct {
	data1 uint32
	data2 uint16
	data3 uint16
	data4 [8]uint8
}

type sSP_DEVICE_INTERFACE_DATA struct {
	cbSize   uint32
	guid     gGUID
	flags    uint32
	reserved uintptr
}

type sSP_DEVICE_INTERFACE_DETAIL_DATA_A struct {
	cbSize     uint32
	devicePath [1]byte
}

type hHIDD_ATTRIBUTES struct {
	size      uint32
	vendorID  uint16
	productID uint16
	version   uint16
}

func listDevices() ([]*Device, error) {
	guid := gGUID{}
	if _, _, err := hidD_GetHidGuid.Call(uintptr(unsafe.Pointer(&guid))); err.(syscall.Errno) != 0 {
		return nil, err
	}

	devInfo, _, err := setupDiGetClassDevsA.Call(uintptr(unsafe.Pointer(&guid)), 0, 0, uintptr(sDIGCF_PRESENT|sDIGCF_DEVICEINTERFACE))
	if err.(syscall.Errno) != 0 {
		return nil, err
	}
	defer setupDiDestroyDeviceInfoList.Call(devInfo)

	idx := uint32(0)
	rv := []*Device{}

	for {
		itf := sSP_DEVICE_INTERFACE_DATA{}
		itf.cbSize = uint32(unsafe.Sizeof(itf))

		b, _, err := setupDiEnumDeviceInterfaces.Call(devInfo, 0, uintptr(unsafe.Pointer(&guid)), uintptr(idx), uintptr(unsafe.Pointer(&itf)))
		idx++
		if b == 0 {
			break
		}
		if err.(syscall.Errno) != 0 {
			continue
		}

		reqSize := uint32(0)
		_, _, err = setupDiGetDeviceInterfaceDetailA.Call(devInfo, uintptr(unsafe.Pointer(&itf)), 0, uintptr(uint32(0)), uintptr(unsafe.Pointer(&reqSize)), 0)
		if err.(syscall.Errno) != syscall.ERROR_INSUFFICIENT_BUFFER {
			continue
		}

		detailBuf := make([]byte, reqSize)
		detail := (*sSP_DEVICE_INTERFACE_DETAIL_DATA_A)(unsafe.Pointer(&detailBuf[0]))
		detail.cbSize = uint32(unsafe.Sizeof(sSP_DEVICE_INTERFACE_DETAIL_DATA_A{}))

		_, _, err = setupDiGetDeviceInterfaceDetailA.Call(devInfo, uintptr(unsafe.Pointer(&itf)), uintptr(unsafe.Pointer(detail)), uintptr(reqSize), 0, 0)
		if err.(syscall.Errno) != 0 {
			continue
		}

		path := strings.TrimSpace(string(detailBuf[unsafe.Offsetof(detail.devicePath) : len(detailBuf)-1]))

		d := func() *Device {
			f, err := os.OpenFile(path, os.O_RDWR, 0755)
			if err != nil {
				return nil
			}
			defer f.Close()

			rv := &Device{
				path: path,
			}

			attr := hHIDD_ATTRIBUTES{}
			_, _, err = hidD_GetAttributes.Call(f.Fd(), uintptr(unsafe.Pointer(&attr)))
			if err.(syscall.Errno) != 0 {
				return nil
			}
			rv.vendorId = attr.vendorID
			rv.productId = attr.productID
			rv.version = attr.version

			buf := make([]uint16, 4092/2)

			_, _, err = hidD_GetManufacturerString.Call(f.Fd(), uintptr(unsafe.Pointer(&buf[0])), unsafe.Sizeof(buf))
			if err.(syscall.Errno) == 0 {
				rv.manufacturer = syscall.UTF16ToString(buf)
			}

			_, _, err = hidD_GetProductString.Call(f.Fd(), uintptr(unsafe.Pointer(&buf[0])), unsafe.Sizeof(buf))
			if err.(syscall.Errno) == 0 {
				rv.product = syscall.UTF16ToString(buf)
			}

			_, _, err = hidD_GetSerialNumberString.Call(f.Fd(), uintptr(unsafe.Pointer(&buf[0])), unsafe.Sizeof(buf))
			if err.(syscall.Errno) == 0 {
				rv.serialNumber = syscall.UTF16ToString(buf)
			}

			return rv
		}()

		if d != nil {
			rv = append(rv, d)
		}
	}

	return rv, nil
}

func (d *Device) lock() error {
	return nil
}
