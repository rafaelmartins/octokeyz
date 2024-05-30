// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: GPL-2.0

#include <avr/eeprom.h>
#include <usbdrv/usbdrv.h>
#include "serialnumber.h"

static uint8_t *sn_addr = (uint8_t*) 1;
static uint8_t rand_bytes[32] __attribute__((section(".noinit")));

int usbDescriptorStringSerialNumber[9];


static inline char
byte2hex(uint8_t v)
{
    return v > 9 ? v - 10 + 'a' : v + '0';
}


void
serialnumber_init(void)
{
    uint32_t sn = eeprom_read_dword((uint32_t*) sn_addr);
    if (sn == 0xffffffff) {
        for (uint8_t i = 0, j = 0; j < 4; j++) {
            for (; i < 32 && (rand_bytes[i] == 0 || rand_bytes[i] == 0xff); i++);
            eeprom_write_byte(sn_addr + j, i < 32 ? rand_bytes[i++] : 0xff);
            eeprom_busy_wait();
        }
        sn = eeprom_read_dword((uint32_t*) sn_addr);
    }

    if (sn != 0xffffffff) {
        usbDescriptorStringSerialNumber[0] = USB_STRING_DESCRIPTOR_HEADER(8);
        for (uint8_t i = 0; i < 8; i++)
            usbDescriptorStringSerialNumber[i % 2 ? i : i + 2] = byte2hex((sn >> (4 * i)) & 0xf);
    }
}
