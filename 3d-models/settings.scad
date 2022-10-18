/*
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

thickness = 2;

pcb_width = 72.39;
pcb_length = 27.94;
pcb_screw_distance_x = 66.04;
pcb_screw_distance_y = 21.59;
pcb_screw_padding = 3.175;
pcb_screw_base_height = 8.5;
pcb_screw_base_d = 5;
pcb_screw_d = 2;
pcb_padding_x = thickness + 0.5;
pcb_padding_y = 5;
pcb_pin0_x = 3.81;
pcb_pin0_y = -0.635;
pcb_pin_distance = 2.54;

pcbs_spacer_height = 4;

width = pcb_width + 2 * pcb_padding_x;
length = pcb_length + 2 * pcb_padding_y + pcb_pin_distance;
height = 2 * thickness + 28;

screw_d = 2;
screw_h = 7;
screw_distance_x = width - 4;
screw_distance_y = length - 4;

cable_d = 4.4;

gap = thickness * 0.2;
