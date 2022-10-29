/*
 * b8: A simple USB keypad with 8 programmable buttons.
 *
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: GPL-2.0
 */

#include <stdbool.h>
#include <avr/io.h>
#include <avr/wdt.h>
#include <avr/interrupt.h>
#include <avr/pgmspace.h>
#include <util/delay.h>
#include "usbdrv/usbdrv.h"

const char usbHidReportDescriptor[] PROGMEM = {
    0x05, 0x0C,  // UsagePage(Consumer[12])
    0x09, 0x01,  // UsageId(Consumer Control[1])
    0xA1, 0x01,  // Collection(Application)
    0x09, 0x03,  //     UsageId(Programmable Buttons[3])
    0xA1, 0x04,  //     Collection(NamedArray)
    0x05, 0x09,  //         UsagePage(Button[9])
    0x19, 0x01,  //         UsageIdMin(Button 1[1])
    0x29, 0x08,  //         UsageIdMax(Button 8[8])
    0x15, 0x00,  //         LogicalMinimum(0)
    0x25, 0x01,  //         LogicalMaximum(1)
    0x95, 0x08,  //         ReportCount(8)
    0x75, 0x01,  //         ReportSize(1)
    0x81, 0x02,  //         Input(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, BitField)
    0xC0,        //     EndCollection()
    0x05, 0x08,  //     UsagePage(LED[8])
    0x09, 0x4B,  //     UsageId(Generic Indicator[75])
    0x15, 0x00,  //     LogicalMinimum(0)
    0x25, 0x01,  //     LogicalMaximum(1)
    0x95, 0x01,  //     ReportCount(1)
    0x75, 0x01,  //     ReportSize(1)
    0x91, 0x02,  //     Output(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0x75, 0x07,  //     ReportSize(7)
    0x91, 0x03,  //     Output(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,        // EndCollection()
};

static volatile uint8_t report;
static volatile bool read_report = false;


usbMsgLen_t
usbFunctionSetup(uchar data[8])
{
    usbRequest_t *rq = (void *) data;
    if ((rq->bmRequestType & USBRQ_TYPE_MASK) == USBRQ_TYPE_CLASS) {
        switch (rq->bRequest) {
        case USBRQ_HID_GET_REPORT:
            usbMsgPtr = (void *) &report;
            return sizeof(report);

        case USBRQ_HID_SET_REPORT:
            read_report = true;
            return USB_NO_MSG;
        }
    }
    return 0;
}


uchar
usbFunctionWrite(uchar *data, uchar len)
{
    if (read_report && len == 1) {
        if (data[0] & (1 << 0))
            PORTD |= (1 << 6);
        else
            PORTD &= ~(1 << 6);
    }
    read_report = false;
    return 1;
}


int
main(void)
{
    PORTB = 0xff;
    DDRD = (1 << 6);

    wdt_enable(WDTO_1S);

    usbInit();
    usbDeviceDisconnect();

    wdt_reset();

    uint8_t i = 0xff;
    while (i--)
        _delay_ms(1);

    usbDeviceConnect();

    sei();

    for (;;) {
        wdt_reset();

        usbPoll();

        if (usbInterruptIsReady()) {
            report = ~PINB;
            usbSetInterrupt((void*) &report, sizeof(report));
        }
    }

    return 0;
}
