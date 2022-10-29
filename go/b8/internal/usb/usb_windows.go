package usb

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	rootUsage           = 0x01 // Consumer Control
	rootUsagePage       = 0x0c // Consumer
	buttonAppUsage      = 0x03 // Programmable Buttons
	buttonAppUsagePage  = 0x0c // Consumer
	buttonNAryUsagePage = 0x09 // Button
	buttonNAryUsageMin  = 0x01 // Button 1
	buttonNAryUsageMax  = 0x08 // Button 8
)

const (
	buttonReportId  = 0
	buttonReportLen = 2
	buttonCaps      = 1
)

const (
	evKey       = 1
	buttonMacro = 0x290
)

const (
	hHIDP_STATUS_SUCCESS                 uintptr = 0x00110000
	hHIDP_STATUS_NULL                    uintptr = 0x80110001
	hHIDP_STATUS_INVALID_PREPARSED_DATA  uintptr = 0xc0110001
	hHIDP_STATUS_INVALID_REPORT_TYPE     uintptr = 0xc0110002
	hHIDP_STATUS_INVALID_REPORT_LENGTH   uintptr = 0xc0110003
	hHIDP_STATUS_USAGE_NOT_FOUND         uintptr = 0xc0110004
	hHIDP_STATUS_VALUE_OUT_OF_RANGE      uintptr = 0xc0110005
	hHIDP_STATUS_BAD_LOG_PHY_VALUES      uintptr = 0xc0110006
	hHIDP_STATUS_BUFFER_TOO_SMALL        uintptr = 0xc0110007
	hHIDP_STATUS_INTERNAL_ERROR          uintptr = 0xc0110008
	hHIDP_STATUS_I8042_TRANS_UNKNOWN     uintptr = 0xc0110009
	hHIDP_STATUS_INCOMPATIBLE_REPORT_ID  uintptr = 0xc011000a
	hHIDP_STATUS_NOT_VALUE_ARRAY         uintptr = 0xc011000b
	hHIDP_STATUS_IS_VALUE_ARRAY          uintptr = 0xc011000c
	hHIDP_STATUS_DATA_INDEX_NOT_FOUND    uintptr = 0xc011000d
	hHIDP_STATUS_DATA_INDEX_OUT_OF_RANGE uintptr = 0xc011000e
	hHIDP_STATUS_BUTTON_NOT_PRESSED      uintptr = 0xc011000f
	hHIDP_STATUS_REPORT_DOES_NOT_EXIST   uintptr = 0xc0110010
	hHIDP_STATUS_NOT_IMPLEMENTED         uintptr = 0xc0110020
)

const (
	hHidP_Input = 0
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
	hidD_FreePreparsedData     = hid.NewProc("HidD_FreePreparsedData")
	hidD_GetAttributes         = hid.NewProc("HidD_GetAttributes")
	hidD_GetHidGuid            = hid.NewProc("HidD_GetHidGuid")
	hidD_GetManufacturerString = hid.NewProc("HidD_GetManufacturerString")
	hidD_GetPreparsedData      = hid.NewProc("HidD_GetPreparsedData")
	hidD_GetProductString      = hid.NewProc("HidD_GetProductString")
	hidD_GetSerialNumberString = hid.NewProc("HidD_GetSerialNumberString")
	hidP_GetButtonCaps         = hid.NewProc("HidP_GetButtonCaps")
	hidP_GetCaps               = hid.NewProc("HidP_GetCaps")
)

var hHIDP_STATUS = map[uintptr]string{
	hHIDP_STATUS_SUCCESS:                 "HIDP_STATUS_SUCCESS",
	hHIDP_STATUS_NULL:                    "HIDP_STATUS_NULL",
	hHIDP_STATUS_INVALID_PREPARSED_DATA:  "HIDP_STATUS_INVALID_PREPARSED_DATA",
	hHIDP_STATUS_INVALID_REPORT_TYPE:     "HIDP_STATUS_INVALID_REPORT_TYPE",
	hHIDP_STATUS_INVALID_REPORT_LENGTH:   "HIDP_STATUS_INVALID_REPORT_LENGTH",
	hHIDP_STATUS_USAGE_NOT_FOUND:         "HIDP_STATUS_USAGE_NOT_FOUND",
	hHIDP_STATUS_VALUE_OUT_OF_RANGE:      "HIDP_STATUS_VALUE_OUT_OF_RANGE",
	hHIDP_STATUS_BAD_LOG_PHY_VALUES:      "HIDP_STATUS_BAD_LOG_PHY_VALUES",
	hHIDP_STATUS_BUFFER_TOO_SMALL:        "HIDP_STATUS_BUFFER_TOO_SMALL",
	hHIDP_STATUS_INTERNAL_ERROR:          "HIDP_STATUS_INTERNAL_ERROR",
	hHIDP_STATUS_I8042_TRANS_UNKNOWN:     "HIDP_STATUS_I8042_TRANS_UNKNOWN",
	hHIDP_STATUS_INCOMPATIBLE_REPORT_ID:  "HIDP_STATUS_INCOMPATIBLE_REPORT_ID",
	hHIDP_STATUS_NOT_VALUE_ARRAY:         "HIDP_STATUS_NOT_VALUE_ARRAY",
	hHIDP_STATUS_IS_VALUE_ARRAY:          "HIDP_STATUS_IS_VALUE_ARRAY",
	hHIDP_STATUS_DATA_INDEX_NOT_FOUND:    "HIDP_STATUS_DATA_INDEX_NOT_FOUND",
	hHIDP_STATUS_DATA_INDEX_OUT_OF_RANGE: "HIDP_STATUS_DATA_INDEX_OUT_OF_RANGE",
	hHIDP_STATUS_BUTTON_NOT_PRESSED:      "HIDP_STATUS_BUTTON_NOT_PRESSED",
	hHIDP_STATUS_REPORT_DOES_NOT_EXIST:   "HIDP_STATUS_REPORT_DOES_NOT_EXIST",
	hHIDP_STATUS_NOT_IMPLEMENTED:         "HIDP_STATUS_NOT_IMPLEMENTED",
}

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

type hHIDP_CAPS struct {
	usage                     uint16
	usagePage                 uint16
	inputReportByteLength     uint16
	outputReportByteLength    uint16
	featureReportByteLength   uint16
	reserved                  [17]uint16
	numberLinkCollectionNodes uint16
	numberInputButtonCaps     uint16
	numberInputValueCaps      uint16
	numberInputDataIndices    uint16
	numberOutputButtonCaps    uint16
	numberOutputValueCaps     uint16
	numberOutputDataIndices   uint16
	numberFeatureButtonCaps   uint16
	numberFeatureValueCaps    uint16
	numberFeatureDataIndices  uint16
}

type hHIDP_BUTTON_CAPS struct {
	usagePage         uint16
	reportID          uint8
	isAlias           bool
	bitField          uint16
	linkCollection    uint16
	linkUsage         uint16
	linkUsagePage     uint16
	isRange           bool
	isStringRange     bool
	isDesignatorRange bool
	isAbsolute        bool
	reportCount       uint16
	reserved2         uint16
	reserved          [9]uint32

	rangeUsageMin      uint16
	rangeUsageMax      uint16
	rangeStringMin     uint16
	rangeStringMax     uint16
	rangeDesignatorMin uint16
	rangeDesignatorMax uint16
	rangeDataIndexMin  uint16
	rangeDataIndexMax  uint16
}

type platformContext struct {
	path       *uint16
	handle     syscall.Handle
	caps       hHIDP_CAPS
	buttonCaps hHIDP_BUTTON_CAPS
	data       uint8
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
		pathp, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			continue
		}

		d := func() *Device {
			h, err := syscall.CreateFile(pathp, syscall.GENERIC_READ|syscall.GENERIC_WRITE, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_EXISTING, 0, -1)
			if err != nil {
				return nil
			}
			defer syscall.CloseHandle(h)

			rv := &Device{
				path: path,
				open: false,
				pctx: platformContext{
					path: pathp,
				},
			}

			attr := hHIDD_ATTRIBUTES{}
			_, _, err = hidD_GetAttributes.Call(uintptr(h), uintptr(unsafe.Pointer(&attr)))
			if err.(syscall.Errno) != 0 {
				return nil
			}
			rv.vendorId = attr.vendorID
			rv.productId = attr.productID

			buf := make([]uint16, 4092/2)

			_, _, err = hidD_GetManufacturerString.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), unsafe.Sizeof(buf))
			if err.(syscall.Errno) == 0 {
				rv.manufacturer = syscall.UTF16ToString(buf)
			}

			_, _, err = hidD_GetProductString.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), unsafe.Sizeof(buf))
			if err.(syscall.Errno) == 0 {
				rv.product = syscall.UTF16ToString(buf)
			}

			_, _, err = hidD_GetSerialNumberString.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), unsafe.Sizeof(buf))
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

func (d *Device) Open() error {
	if d.open {
		return errors.New("usb: device is open")
	}

	h, err := syscall.CreateFile(d.pctx.path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_EXISTING, 0, -1)
	if err != nil {
		return err
	}
	d.open = true
	d.pctx.handle = h

	var preparsed uintptr

	b, _, err := hidD_GetPreparsedData.Call(uintptr(d.pctx.handle), uintptr(unsafe.Pointer(&preparsed)))
	if err.(syscall.Errno) != 0 {
		d.Close()
		return err
	}
	if b == 0 {
		d.Close()
		return errors.New("usb: failed to preparse device data")
	}
	defer hidD_FreePreparsedData.Call(preparsed)

	status, _, err := hidP_GetCaps.Call(preparsed, uintptr(unsafe.Pointer(&d.pctx.caps)))
	if err.(syscall.Errno) != 0 {
		d.Close()
		return err
	}
	if status != hHIDP_STATUS_SUCCESS {
		d.Close()
		return fmt.Errorf("usb: NTSTATUS = %s", hHIDP_STATUS[status])
	}

	if d.pctx.caps.usage != rootUsage {
		d.Close()
		return fmt.Errorf("usb: device reported wrong usage: %d", d.pctx.caps.usage)
	}
	if d.pctx.caps.usagePage != rootUsagePage {
		d.Close()
		return fmt.Errorf("usb: device reported wrong usage page: %d", d.pctx.caps.usagePage)
	}
	if d.pctx.caps.inputReportByteLength != buttonReportLen {
		d.Close()
		return fmt.Errorf("usb: device reported wrong input report byte length: %d", d.pctx.caps.inputReportByteLength)
	}
	if d.pctx.caps.numberInputButtonCaps != buttonCaps {
		d.Close()
		return fmt.Errorf("usb: device reported wrong number of button capabilities: %d", d.pctx.caps.numberInputButtonCaps)
	}
	if d.pctx.caps.numberInputDataIndices != buttonNAryUsageMax-buttonNAryUsageMin+1 {
		d.Close()
		return fmt.Errorf("usb: device reported wrong number of input data indices: %d", d.pctx.caps.numberInputDataIndices)
	}

	status, _, err = hidP_GetButtonCaps.Call(uintptr(hHidP_Input), uintptr(unsafe.Pointer(&d.pctx.buttonCaps)), uintptr(unsafe.Pointer(&d.pctx.caps.numberInputButtonCaps)), preparsed)
	if err.(syscall.Errno) != 0 {
		d.Close()
		return err
	}
	if status != hHIDP_STATUS_SUCCESS {
		d.Close()
		return fmt.Errorf("usb: NTSTATUS = %s", hHIDP_STATUS[status])
	}

	if d.pctx.buttonCaps.usagePage != buttonNAryUsagePage {
		d.Close()
		return fmt.Errorf("usb: device reported wrong button usage: %d", d.pctx.buttonCaps.usagePage)
	}
	if d.pctx.buttonCaps.reportID != buttonReportId {
		d.Close()
		return fmt.Errorf("usb: device reported wrong button report id: %d", d.pctx.buttonCaps.reportID)
	}
	if d.pctx.buttonCaps.linkUsage != buttonAppUsage {
		d.Close()
		return fmt.Errorf("usb: device reported wrong button application usage: %d", d.pctx.buttonCaps.linkUsage)
	}
	if d.pctx.buttonCaps.linkUsagePage != rootUsagePage {
		d.Close()
		return fmt.Errorf("usb: device reported wrong usage: %d", d.pctx.buttonCaps.linkUsagePage)
	}
	if !d.pctx.buttonCaps.isRange {
		d.Close()
		return errors.New("usb: device reported that usages are not organized in range")
	}
	if d.pctx.buttonCaps.rangeUsageMin != buttonNAryUsageMin {
		d.Close()
		return fmt.Errorf("usb: device reported wrong button usage minimum: %d", d.pctx.buttonCaps.rangeUsageMin)
	}
	if d.pctx.buttonCaps.rangeUsageMax != buttonNAryUsageMax {
		d.Close()
		return fmt.Errorf("usb: device reported wrong button usage maximum: %d", d.pctx.buttonCaps.rangeUsageMax)
	}

	return nil
}

func (d *Device) Close() error {
	if !d.open {
		return nil
	}

	if err := syscall.CloseHandle(d.pctx.handle); err != nil {
		return err
	}

	d.open = false
	d.pctx.handle = syscall.InvalidHandle

	return nil
}

func (d *Device) Read() ([]*Event, error) {
	if !d.open {
		return nil, errors.New("usb: device is not open")
	}

	buf := make([]byte, 64*buttonReportLen)
	n := uint32(0)
	if err := syscall.ReadFile(d.pctx.handle, buf, &n, nil); err != nil {
		return nil, err
	}
	if n%buttonReportLen != 0 {
		return nil, errors.New("usb: failed to read hid report")
	}

	rv := []*Event{}

	for i := 0; i < int(n); i += buttonReportLen {
		if buf[i] != buttonReportId {
			continue
		}

		if buf[i+1] == d.pctx.data {
			continue
		}

		t := time.Now()
		for j := d.pctx.buttonCaps.rangeDataIndexMin; j <= d.pctx.buttonCaps.rangeDataIndexMax; j++ {
			if (d.pctx.data & (1 << j)) != (buf[i+1] & (1 << j)) {
				rv = append(rv, &Event{
					Time:  t,
					Type:  evKey,
					Code:  buttonMacro + j,
					Value: int32(buf[i+1] & (1 << j)),
				})
			}
		}
		d.pctx.data = buf[i+1]
	}

	return rv, nil
}
