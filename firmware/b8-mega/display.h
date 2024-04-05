/*
 * b8: A simple USB keypad with 8 programmable buttons.
 *
 * SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: BSD-3-Clause
 */

#pragma once

#include <stdbool.h>
#include <stdint.h>

typedef enum {
    DISPLAY_HALIGN_LEFT = 1,
    DISPLAY_HALIGN_RIGHT,
    DISPLAY_HALIGN_CENTER,
} display_halign_t;

void display_init(void);
bool display_line(uint8_t line, const char *str, display_halign_t align);
bool display_line_from_report(uint8_t *buf, uint8_t len);
void display_clear(void);
void display_task(void);
