/*
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

include <settings.scad>

width_ = 3;
length_ = 2 * led_margin + led_length;
height_ = pcb_screw_base_height - 2.5;

cube([led_margin, width_, height_]);
translate([led_margin, (width_ - led_width) / 2, 0])
    cube([led_length, led_width, height_]);
translate([length_ - led_margin, 0, 0])
    cube([led_margin, width_, height_]);
