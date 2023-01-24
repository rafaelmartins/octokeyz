/*
 * b8: A simple USB keypad with 8 programmable buttons.
 *
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: GPL-2.0
 */

#include <avr/eeprom.h>
#include <usbdrv/usbdrv.h>
#include "serialnumber.h"

static uint8_t *sn_addr = (uint8_t*) 1;
static uint8_t rand_bytes[32] __attribute__((section(".noinit")));

int usbDescriptorStringSerialNumber[9];


static char
byte2hex(uint8_t v)
{
    return v > 9 ? v - 10 + 'a' : v + '0';
}


void
serialnumber_init(void)
{
    uint32_t sn = eeprom_read_dword((uint32_t*) sn_addr);
    if (sn == 0 || sn == 0xffffffff) {
        uint8_t j = 0;
        for (uint8_t i = 0; i < 32 && j < 4; i++) {
            if (rand_bytes[i] != 0 && rand_bytes[i] != 0xff) {
                eeprom_write_byte(sn_addr + j++, rand_bytes[i]);
            }
        }
        while (j < 4) {
            eeprom_write_byte(sn_addr + j++, 0);
        }
        sn = eeprom_read_dword((uint32_t*) sn_addr);
    }

    usbDescriptorStringSerialNumber[0] = USB_STRING_DESCRIPTOR_HEADER(8);
    for (uint8_t i = 0; i < 8; i++)
        usbDescriptorStringSerialNumber[i % 2 ? i : i + 2] = byte2hex((sn >> (4 * i)) & 0xf);
}
