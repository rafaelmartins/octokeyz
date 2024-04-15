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
#include <usbdrv/usbdrv.h>
#include "bits.h"

#if !(defined(__AVR_ATtiny2313__) || defined(__AVR_ATtiny2313A__))
#include "serialnumber.h"
#endif

FUSES = {
    .low = FUSE_SUT1,
    .high = FUSE_SPIEN,
    .extended = EFUSE_DEFAULT,
};
LOCKBITS = LOCKBITS_DEFAULT;

const char usbHidReportDescriptor[] PROGMEM = {
    0x05, 0x0C,    // UsagePage(Consumer[0x000C])
    0x09, 0x01,    // UsageId(Consumer Control[0x0001])
    0xA1, 0x01,    // Collection(Application)
    0x85, 0x01,    //     ReportId(1)
    0x09, 0x03,    //     UsageId(Programmable Buttons[0x0003])
    0xA1, 0x02,    //     Collection(Logical)
    0x05, 0x09,    //         UsagePage(Button[0x0009])
    0x19, 0x01,    //         UsageIdMin(Button 1[0x0001])
    0x29, 0x08,    //         UsageIdMax(Button 8[0x0008])
    0x15, 0x00,    //         LogicalMinimum(0)
    0x25, 0x01,    //         LogicalMaximum(1)
    0x95, 0x08,    //         ReportCount(8)
    0x75, 0x01,    //         ReportSize(1)
    0x81, 0x02,    //         Input(Data, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, BitField)
    0xC0,          //     EndCollection()
    0x05, 0x08,    //     UsagePage(LED[0x0008])
    0x09, 0x3C,    //     UsageId(Usage Multi Mode Indicator[0x003C])
    0xA1, 0x02,    //     Collection(Logical)
    0x19, 0x3D,    //         UsageIdMin(Indicator On[0x003D])
    0x29, 0x41,    //         UsageIdMax(Indicator Off[0x0041])
    0x15, 0x01,    //         LogicalMinimum(1)
    0x25, 0x05,    //         LogicalMaximum(5)
    0x95, 0x01,    //         ReportCount(1)
    0x75, 0x03,    //         ReportSize(3)
    0x91, 0x00,    //         Output(Data, Array, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,          //     EndCollection()
    0x75, 0x05,    //     ReportSize(5)
    0x91, 0x03,    //     Output(Constant, Variable, Absolute, NoWrap, Linear, PreferredState, NoNullPosition, NonVolatile, BitField)
    0xC0,          // EndCollection()
};

static volatile uint8_t report[2] = {1};
static volatile enum {
    LED_ON = 1,
    LED_FLASH,
    LED_SLOW_BLINK,
    LED_FAST_BLINK,
    LED_OFF,
} led_state = LED_OFF;


usbMsgLen_t
usbFunctionSetup(uchar data[8])
{
    usbRequest_t *rq = (void *) data;
    if ((rq->bmRequestType & USBRQ_TYPE_MASK) == USBRQ_TYPE_CLASS) {
        switch (rq->bRequest) {
        case USBRQ_HID_GET_REPORT:
            report[1] = ~PINB;
            usbMsgPtr = (void *) report;
            return sizeof(report);

        case USBRQ_HID_SET_REPORT:
            return USB_NO_MSG;
        }
    }
    return 0;
}


ISR(TIMER1_COMPA_vect)
{
    switch (led_state) {
    case LED_FLASH:
        led_state = LED_OFF;
        PORT_CLEAR(P_LED);
        TCCR1B = 0;
        break;

    case LED_SLOW_BLINK:
    case LED_FAST_BLINK:
        PORT_FLIP(P_LED);
        break;

    case LED_OFF:
    case LED_ON:
        break;
    }
}


uchar
usbFunctionWrite(uchar *data, uchar len)
{
    if (len == 2 && data[0] == 1) {
        led_state = data[1] & 0b11111;

        TCCR1B = 0;
        TCNT1 = 0;
        OCR1A = 0;

        switch (led_state) {
        case LED_OFF:
            PORT_CLEAR(P_LED);
            break;

        case LED_SLOW_BLINK:
            OCR1A += 2928;  // ~150ms @ F_CPU/1024
            // fallthrough

        case LED_FAST_BLINK:
            OCR1A += 976;   // ~50ms @ F_CPU/1024
            // fallthrough

        case LED_FLASH:
            OCR1A += 976;   // ~50ms @ F_CPU/1024
            TCCR1B = (1 << WGM12) | (1 << CS12) | (1 << CS10);
            // fallthrough

        case LED_ON:
            PORT_SET(P_LED);
            break;

        default:
            break;
        }
    }
    return 1;
}


int
main(void)
{
#if !(defined(__AVR_ATtiny2313__) || defined(__AVR_ATtiny2313A__))
    serialnumber_init();
#endif

    PORTB = 0xff;
    DDR_SET(P_LED);
    PORT_SET(P_LED);

    TCCR1A = 0;
    TIMSK = (1 << OCIE1A);

    wdt_enable(WDTO_2S);

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
            report[1] = ~PINB;
            usbSetInterrupt((void*) report, sizeof(report));
        }
    }

    return 0;
}
