/*
 * SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

include <../lib/ssd1306.scad>

screw_d = 2;
screw_h = 8;
thickness = 2;
gap = 0;

front_screw_padding = 3.5;

pcb_base_d = 5;
pcb_base_screw_d = 1.8;
pcb_screw_padding = 2.54;
pcb_thickness = 1.6;

pcb_padding_x = thickness + 1;
pcb_padding_y = thickness + 1;

display_pcb_length = 39.624;
display_pcb_width = 33.782;

display_pcb0_x = pcb_padding_x;
display_pcb0_y = pcb_padding_y;

display_length = display_pcb_length + 2 * pcb_padding_x;
display_width = display_pcb_width + pcb_padding_y;

switches_pcb_length = 84.582;
switches_pcb_width = 45.974;

switches_length = switches_pcb_length + 2 * pcb_padding_x;
switches_width = switches_pcb_width + 2 * pcb_padding_y;

switches_pcb_base_distance_x = 39.751;
switches_pcb_base_distance_y = 20.447;

oled_socket_height = 8.5;
oled_pin0_y = 5.334;
oled0_x = (display_length - ssd1306_pcb_base_spacing_x) / 2;
oled0_y = display_pcb0_y + oled_pin0_y + ssd1306_pin_distance_y;

display_pcb_base_distance = 34.544;
display_pcb_base_height = 2.8 + oled_socket_height + ssd1306_pcb_thickness + ssd1306_pcb_base_height;

switches_pcb_base_height = 6 - thickness;

usb_height = 4.4;
usb_width = 8.2;
usb_distance_x = pcb_padding_x + 29.210 - usb_width / 2;
usb_distance_z = thickness + display_pcb_base_height + pcb_thickness - 0.1;

display_height = usb_distance_z + usb_height + 0.8 + thickness;
display_height_in = display_pcb_base_height - switches_pcb_base_height - thickness;

switches_height = display_height - display_pcb_base_height + switches_pcb_base_height;

led0_x = 37.211;
led0_y = 15.24;
led_d = 3.2;

key_dim = 16;
key_distance = 19.304;
key0_x = 13.335 - key_dim / 2;
key0_y = key0_x;
