// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#pragma once

#include <stdbool.h>
#include <stdint.h>

#include "display-font.h"

#define display_address 0x3c
#define display_screen_width 128
#define display_screen_height 64
#define display_chars_per_line (display_screen_width / (display_font_width + 1))
#define display_lines (display_screen_height / 8)  // info from ssd1306 datasheet

typedef enum {
    DISPLAY_HALIGN_LEFT = 1,
    DISPLAY_HALIGN_RIGHT,
    DISPLAY_HALIGN_CENTER,
} display_halign_t;

bool display_init(void);
bool display_line(uint8_t line, const char *str, display_halign_t align);
bool display_line_from_report(uint8_t *buf, uint8_t len);
void display_clear_line(uint8_t line);
void display_clear(void);
void display_task(void);
