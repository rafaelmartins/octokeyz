// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: CERN-OHL-S-2.0

include <../lib/screw-base.scad>
include <settings.scad>

keycap_model = "B32-13XX";  // supports "B32-13XX" (square) and "B32-16XX" (circle), both from OSRAM

difference() {
    union() {
        cube([width, length, thickness]);
        cube([width, thickness, height]);
        translate([0, length - thickness, 0])
            cube([width, thickness, height]);

        for (i=[0:1]) {
            translate([pcb_padding_x + i * pcb_width, pcb_padding_y + i * pcb_length, thickness]) {
                rotate([0, 0, i * 180]) {
                    translate([pcb_screw1_x, pcb_screw1_y, 0])
                        cylinder(pcb_screw_base_height, d=pcb_screw_base_d, $fn=20);
                    translate([pcb_screw2_x, pcb_screw2_y, 0])
                        cylinder(pcb_screw_base_height, d=pcb_screw_base_d, $fn=20);
                    translate([pcb_screw3_x, pcb_screw3_y, 0])
                        cylinder(pcb_screw_base_height, d=pcb_screw_base_d, $fn=20);
                }
            }
        }
    }

    translate([pcb_padding_x, pcb_padding_y, 0]) {
        for (i=[0:1]) {
            translate([i * pcb_width, i * pcb_length, thickness]) {
                rotate([0, 0, i * 180]) {
                    translate([pcb_screw1_x, pcb_screw1_y, 0])
                        cylinder(pcb_screw_base_height, d=pcb_screw_d - 0.1, $fn=20);
                    translate([pcb_screw2_x, pcb_screw2_y, 0])
                        cylinder(pcb_screw_base_height, d=pcb_screw_d - 0.1, $fn=20);
                    translate([pcb_screw3_x, pcb_screw3_y, 0])
                        cylinder(pcb_screw_base_height, d=pcb_screw_d - 0.1, $fn=20);
                }
            }
        }

        for (i=[0:3])
            for (j=[0:1])
                if (keycap_model == "B32-13XX")
                    translate([button0_x - button_dim / 2 + i * button_distance,
                               button0_y - button_dim / 2 + j * button_distance, 0])
                        cube([button_dim, button_dim, thickness]);
                else
                    translate([button0_x + i * button_distance, button0_y + j * button_distance, 0])
                        cylinder(thickness, d=button_d, $fn=20);

        translate([led_x, led_y, 0])
            cylinder(thickness, d=led_d, $fn=20);
    }

    translate([usb_distance_x, 0, usb_distance_z])
        cube([usb_width, thickness, height - usb_distance_z]);

    for (i=[0:1])
        for (j=[0:1])
            translate([screw_base_dim(screw_d) / 2 + i * (width - screw_base_dim(screw_d)),
                       j * (length - thickness),
                       height - screw_base_dim(screw_d) / 2])
                rotate([-90, 0, 0])
                    cylinder(thickness, d=screw_d * 1.1, $fn=20);

    translate([(thickness - gap) / 2, thickness, thickness / 2])
        cube([thickness / 2 + gap, length - 2 * thickness, thickness / 2]);
    translate([width - thickness - gap / 2, thickness, thickness / 2])
        cube([thickness / 2 + gap, length - 2 * thickness, thickness / 2]);
}
