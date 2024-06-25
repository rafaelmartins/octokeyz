// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#include <stdlib.h>

#include <stm32f0xx.h>

#include <usbd.h>
#include <usb-std-hid.h>

#include "bootloader.h"
#include "display.h"
#include "led.h"

#define BOOTLOADER_COMBO (GPIO_IDR_0 | GPIO_IDR_4)

static bool display_available = false;


void
clock_init(void)
{
    // 1 flash wait cycle required to operate @ 48MHz (RM0091 section 3.5.1)
    FLASH->ACR &= ~FLASH_ACR_LATENCY;
    FLASH->ACR |= FLASH_ACR_LATENCY;
    while ((FLASH->ACR & FLASH_ACR_LATENCY) != FLASH_ACR_LATENCY);

    RCC->CR2 |= RCC_CR2_HSI48ON;
    while ((RCC->CR2 & RCC_CR2_HSI48RDY) != RCC_CR2_HSI48RDY);

    RCC->CFGR &= ~(RCC_CFGR_HPRE | RCC_CFGR_PPRE | RCC_CFGR_SW);
    RCC->CFGR |= RCC_CFGR_HPRE_DIV1 | RCC_CFGR_PPRE_DIV1 | RCC_CFGR_SW_HSI48;
    while((RCC->CFGR & RCC_CFGR_SWS) != RCC_CFGR_SWS_HSI48);

    SystemCoreClock = 48000000;
}


void
usbd_in_cb(uint8_t ept)
{
    if (ept != 1)
        return;

    uint8_t v[] = {
        1,
        ~((uint8_t) (GPIOA->IDR)),
    };
    usbd_in(ept, &v, sizeof(v));
}


void
usbd_out_cb(uint8_t ept)
{
    if (ept != 1)
        return;

    uint8_t buf[USBD_EP1_OUT_SIZE];
    uint16_t len = usbd_out(ept, buf, sizeof(buf));

    switch (buf[0]) {
    case 1:
        if (len == 2)
            led_set_state(buf[1] & 0b11111);
        break;

    case 2:
        if (len > 1)
            display_line_from_report(buf + 1, len - 1);
        break;
    }
}


bool
usbd_ctrl_request_handle_class_cb(usb_ctrl_request_t *req)
{
    switch (req->bRequest) {
    case USB_REQ_HID_GET_REPORT:
        if ((req->wValue >> 8) != 3)
            break;

        switch ((uint8_t) (req->wValue)) {
        case 1:
            {
                uint8_t data[] = {
                    1,
                    display_available ? (1 << 0) : 0,
                };
                usbd_control_in(data, sizeof(data), req->wLength);
                return true;
            }
            break;

        case 2:
            {
                uint8_t data[] = {
                    2,
                    display_available ? display_lines : 0,
                    display_available ? display_chars_per_line: 0,
                };
                usbd_control_in(data, sizeof(data), req->wLength);
                return true;
            }
            break;
        }
        break;
    }
    return false;
}


void
usbd_reset_hook_cb(bool before)
{
    if (before) {
        led_set_state(LED_ON);
        display_clear();
        display_line(1, "octokeyz", DISPLAY_HALIGN_CENTER);
        display_line(3, OCTOKEYZ_VERSION, DISPLAY_HALIGN_CENTER);
        display_line(6, "Waiting for USB ...", DISPLAY_HALIGN_CENTER);
    }
}


void
usbd_set_address_hook_cb(uint8_t addr)
{
    (void) addr;

    led_set_state(LED_OFF);
    display_clear();
}


int
main(void)
{
    bootloader_entry();

    RCC->AHBENR |= RCC_AHBENR_GPIOAEN;

    GPIOA->PUPDR |=
        GPIO_PUPDR_PUPDR0_0 | GPIO_PUPDR_PUPDR1_0 | GPIO_PUPDR_PUPDR2_0 |
        GPIO_PUPDR_PUPDR3_0 | GPIO_PUPDR_PUPDR4_0 | GPIO_PUPDR_PUPDR5_0 |
        GPIO_PUPDR_PUPDR6_0 | GPIO_PUPDR_PUPDR7_0;

    // wait a little bit until the pull-ups are stable.
    for (__IO uint16_t i = 0xffff; i; i--);

    if ((uint8_t) GPIOA->IDR == (uint8_t) ~BOOTLOADER_COMBO)
        bootloader_reset();

    clock_init();
    led_init();
    display_available = display_init();

    usbd_init();

    while (true) {
        usbd_task();
        display_task();
    }

    return 0;
}
