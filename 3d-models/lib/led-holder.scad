/*
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

module led_holder(height) {
    led_width = 1.8;
    led_length = 1;
    led_margin = 1;

    width_ = 3;
    length_ = 2 * led_margin + led_length;

    cube([led_margin, width_, height]);
    translate([led_margin, (width_ - led_width) / 2, 0])
        cube([led_length, led_width, height]);
    translate([length_ - led_margin, 0, 0])
        cube([led_margin, width_, height]);
}
