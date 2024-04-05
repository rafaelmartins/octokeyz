/*
 * b8: A simple USB keypad with 8 programmable buttons.
 *
 * SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: BSD-3-Clause
 */

#pragma once

typedef enum {
    LED_ON = 1,
    LED_FLASH,
    LED_SLOW_BLINK,
    LED_FAST_BLINK,
    LED_OFF,
} led_state_t;

void led_init(void);
void led_set_state(led_state_t state);
