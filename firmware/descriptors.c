// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <stdlib.h>

#include <usbd.h>
#include <usb-std-hid.h>

#define MANUFACTURER u"rgm.io"
#define PRODUCT      u"octokeyz"

#define ID_VENDOR  0x1d50
#define ID_PRODUCT 0x6184

// +----------+---------+-------------------+
// | ReportId | Kind    | ReportSizeInBytes |
// +----------+---------+-------------------+
// |        1 | Input   |                 1 |
// +----------+---------+-------------------+
// |        1 | Output  |                 1 |
// +----------+---------+-------------------+
// |        1 | Feature |                 1 |
// +----------+---------+-------------------+
// |        2 | Output  |                23 |
// +----------+---------+-------------------+
// |        2 | Feature |                 3 |
// +----------+---------+-------------------+
// |        3 | Output  |                 2 |
// +----------+---------+-------------------+
static const uint8_t hid_report_descriptor[] = {
    0x06, 0x00, 0xFF,                // UsagePage(octokeyz[0xFF00])
    0x09, 0x01,                      // UsageId(octokeyz[0x0001])
    0xA1, 0x01,                      // Collection(Application)
    0x85, 0x01,                      //     ReportId(1)
    0x09, 0x03,                      //     UsageId(Keys[0x0003])
    0xA1, 0x02,                      //     Collection(Logical)
    0x06, 0x01, 0xFF,                //         UsagePage(octokeyz Key[0xFF01])
    0x19, 0x01,                      //         UsageIdMin(Key 1[0x0001])
    0x29, 0x08,                      //         UsageIdMax(Key 8[0x0008])
    0x15, 0x00,                      //         LogicalMinimum(0)
    0x25, 0x01,                      //         LogicalMaximum(1)
    0x95, 0x08,                      //         ReportCount(8)
    0x75, 0x01,                      //         ReportSize(1)
    0x81, 0x02,                      //         Input(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, BitField)
    0xC0,                            //     EndCollection()
    0x06, 0x02, 0xFF,                //     UsagePage(octokeyz LED[0xFF02])
    0x09, 0x01,                      //     UsageId(LED[0x0001])
    0xA1, 0x02,                      //     Collection(Logical)
    0x19, 0x10,                      //         UsageIdMin(LED On[0x0010])
    0x29, 0x14,                      //         UsageIdMax(LED Off[0x0014])
    0x15, 0x01,                      //         LogicalMinimum(1)
    0x25, 0x05,                      //         LogicalMaximum(5)
    0x95, 0x01,                      //         ReportCount(1)
    0x75, 0x03,                      //         ReportSize(3)
    0x91, 0x00,                      //         Output(Data, Array, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,                            //     EndCollection()
    0x75, 0x05,                      //     ReportSize(5)
    0x91, 0x03,                      //     Output(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x06, 0x00, 0xFF,                //     UsagePage(octokeyz[0xFF00])
    0x09, 0x02,                      //     UsageId(Capabilities[0x0002])
    0xA1, 0x02,                      //     Collection(Logical)
    0x09, 0x10,                      //         UsageId(With Display[0x0010])
    0x15, 0x00,                      //         LogicalMinimum(0)
    0x25, 0x01,                      //         LogicalMaximum(1)
    0x75, 0x01,                      //         ReportSize(1)
    0xB1, 0x03,                      //         Feature(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,                            //     EndCollection()
    0x75, 0x07,                      //     ReportSize(7)
    0xB1, 0x03,                      //     Feature(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x85, 0x02,                      //     ReportId(2)
    0x06, 0x03, 0xFF,                //     UsagePage(octokeyz Display[0xFF03])
    0x09, 0x01,                      //     UsageId(Display Capabilities[0x0001])
    0xA1, 0x02,                      //     Collection(Logical)
    0x09, 0x10,                      //         UsageId(Display Lines[0x0010])
    0x09, 0x11,                      //         UsageId(Display Characters per Line[0x0011])
    0x26, 0xFF, 0x00,                //         LogicalMaximum(255)
    0x95, 0x02,                      //         ReportCount(2)
    0x75, 0x08,                      //         ReportSize(8)
    0xB1, 0x03,                      //         Feature(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x09, 0x12,                      //         UsageId(With Display Clear[0x0012])
    0x25, 0x01,                      //         LogicalMaximum(1)
    0x95, 0x01,                      //         ReportCount(1)
    0x75, 0x01,                      //         ReportSize(1)
    0xB1, 0x03,                      //         Feature(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,                            //     EndCollection()
    0x75, 0x07,                      //     ReportSize(7)
    0xB1, 0x03,                      //     Feature(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x09, 0x02,                      //     UsageId(Display Data[0x0002])
    0xA1, 0x02,                      //     Collection(Logical)
    0x09, 0x20,                      //         UsageId(Line[0x0020])
    0x25, 0x1F,                      //         LogicalMaximum(31)
    0x75, 0x05,                      //         ReportSize(5)
    0x91, 0x02,                      //         Output(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x75, 0x03,                      //         ReportSize(3)
    0x91, 0x03,                      //         Output(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x19, 0x30,                      //         UsageIdMin(Alignment Left[0x0030])
    0x29, 0x32,                      //         UsageIdMax(Alignment Center[0x0032])
    0x15, 0x01,                      //         LogicalMinimum(1)
    0x25, 0x03,                      //         LogicalMaximum(3)
    0x75, 0x02,                      //         ReportSize(2)
    0x91, 0x00,                      //         Output(Data, Array, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x75, 0x06,                      //         ReportSize(6)
    0x91, 0x03,                      //         Output(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x09, 0x21,                      //         UsageId(Line Data[0x0021])
    0x15, 0x00,                      //         LogicalMinimum(0)
    0x26, 0xFF, 0x00,                //         LogicalMaximum(255)
    0x95, 0x15,                      //         ReportCount(21)
    0x75, 0x08,                      //         ReportSize(8)
    0x91, 0x02,                      //         Output(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,                            //     EndCollection()
    0x85, 0x03,                      //     ReportId(3)
    0x09, 0x04,                      //     UsageId(Display Clear[0x0004])
    0xA1, 0x02,                      //     Collection(Logical)
    0x09, 0x40,                      //         UsageId(Delay Miliseconds[0x0040])
    0x27, 0xFF, 0xFF, 0x00, 0x00,    //         LogicalMaximum(65,535)
    0x95, 0x01,                      //         ReportCount(1)
    0x75, 0x10,                      //         ReportSize(16)
    0x91, 0x02,                      //         Output(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,                            //     EndCollection()
    0xC0,                            // EndCollection()
};

static const usb_device_descriptor_t device_descriptor = {
    .bLength = sizeof(usb_device_descriptor_t),
    .bDescriptorType = USB_DESCR_TYPE_DEVICE,
    .bcdUSB = 0x0200,
    .bDeviceClass = 0,
    .bDeviceSubClass = 0,
    .bDeviceProtocol = 0,
    .bMaxPacketSize0 = USBD_EP0_SIZE,
    .idVendor = ID_VENDOR,
    .idProduct = ID_PRODUCT,
    .bcdDevice = 0x0101,
    .iManufacturer = 1,
    .iProduct = 2,
    .iSerialNumber = 3,
    .bNumConfigurations = 1,
};

const usb_device_descriptor_t*
usbd_get_device_descriptor_cb(void)
{
    return &device_descriptor;
}


typedef __PACKED_STRUCT {
    usb_config_descriptor_t config_descriptor;
    usb_interface_descriptor_t interface_descriptor;
    usb_hid_descriptor_t hid_descriptor;
    usb_endpoint_descriptor_t endpoint_in_descriptor;
    usb_endpoint_descriptor_t endpoint_out_descriptor;
} config_descriptor_t;

static const config_descriptor_t config_descriptor = {
    .config_descriptor = {
        .bLength = sizeof(usb_config_descriptor_t),
        .bDescriptorType = USB_DESCR_TYPE_CONFIGURATION,
        .wTotalLength = sizeof(config_descriptor_t),
        .bNumInterfaces = 1,
        .bConfigurationValue = 1,
        .iConfiguration = 0,
        .bmAttributes = USB_DESCR_CONFIG_ATTR_RESERVED,
        .bMaxPower = 50, // 100mA
    },
    .interface_descriptor = {
        .bLength = sizeof(usb_interface_descriptor_t),
        .bDescriptorType = USB_DESCR_TYPE_INTERFACE,
        .bInterfaceNumber = 0,
        .bAlternateSetting = 0,
        .bNumEndpoints = 2,
        .bInterfaceClass = USB_DESCR_DEV_CLASS_HID,
        .bInterfaceSubClass = 0,
        .bInterfaceProtocol = 0,
        .iInterface = 0,
    },
    .hid_descriptor = {
        .bLength = sizeof(usb_hid_descriptor_t),
        .bDescriptorType = USB_DESCR_TYPE_HID,
        .bcdHID = 0x0111,
        .bCountryCode = 0,
        .bNumDescriptors = 1,
        .bDescriptorType2 = USB_DESCR_TYPE_HID_REPORT,
        .wDescriptorLength = sizeof(hid_report_descriptor),
    },
    .endpoint_in_descriptor = {
        .bLength = sizeof(usb_endpoint_descriptor_t),
        .bDescriptorType = USB_DESCR_TYPE_ENDPOINT,
        .bEndpointAddress = USB_DESCR_EPT_ADDR_DIR_IN | 1,
        .bmAttributes = USB_DESCR_EPT_ATTR_INTERRUPT,
        .wMaxPacketSize = USBD_EP1_IN_SIZE,
        .bInterval = 10,
        .bRefresh = 0,
        .bSynchAddress = 0,
    },
    .endpoint_out_descriptor = {
        .bLength = sizeof(usb_endpoint_descriptor_t),
        .bDescriptorType = USB_DESCR_TYPE_ENDPOINT,
        .bEndpointAddress = USB_DESCR_EPT_ADDR_DIR_OUT | 1,
        .bmAttributes = USB_DESCR_EPT_ATTR_INTERRUPT,
        .wMaxPacketSize = USBD_EP1_OUT_SIZE,
        .bInterval = 10,
        .bRefresh = 0,
        .bSynchAddress = 0,
    },
};

const usb_config_descriptor_t*
usbd_get_config_descriptor_cb(void)
{
    return &config_descriptor.config_descriptor;
}

const usb_interface_descriptor_t*
usbd_get_interface_descriptor_cb(uint16_t itf)
{
    switch (itf) {
    case 0:
        return &config_descriptor.interface_descriptor;
    }
    return NULL;
}


static const usb_string_descriptor_t language = {
    .bLength = 4,
    .bDescriptorType = USB_DESCR_TYPE_STRING,
    .wData = {0x0409},
};

static const usb_string_descriptor_t manufacturer = {
    .bLength = sizeof(MANUFACTURER),
    .bDescriptorType = USB_DESCR_TYPE_STRING,
    .wData = MANUFACTURER,
};

static const usb_string_descriptor_t product = {
    .bLength = sizeof(PRODUCT),
    .bDescriptorType = USB_DESCR_TYPE_STRING,
    .wData = PRODUCT,
};

const usb_string_descriptor_t*
usbd_get_string_descriptor_cb(uint16_t lang, uint8_t idx)
{
    (void) lang;
    switch (idx) {
    case 0:
        return &language;
    case 1:
        return &manufacturer;
    case 2:
        return &product;
    case 3:
        return usbd_serial_internal_string_descriptor();
    }
    return NULL;
}


bool
usbd_ctrl_request_get_descriptor_interface_cb(usb_ctrl_request_t *req)
{
    if (((uint8_t) req->wIndex != 0))
        return false;

    switch (req->bRequest) {
    case USB_REQ_GET_DESCRIPTOR:
        switch (req->wValue >> 8) {
        case USB_DESCR_TYPE_HID:
            usbd_control_in(&config_descriptor.hid_descriptor, config_descriptor.hid_descriptor.bLength, req->wLength);
            return true;

        case USB_DESCR_TYPE_HID_REPORT:
            usbd_control_in(&hid_report_descriptor, sizeof(hid_report_descriptor), req->wLength);
            return true;
        }
        break;
    }
    return false;
}
