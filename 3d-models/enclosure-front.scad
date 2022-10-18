/*
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

include <lib/screw-base.scad>
include <settings.scad>

button0_x = 2.5 * pcb_pin_distance;
button0_y = 2 * pcb_pin_distance;
button_distance = 6 * pcb_pin_distance;


difference() {
    union() {
        cube([width, length, thickness]);
        cube([thickness, length, height]);
        translate([width - thickness, 0, 0])
            cube([thickness, length, height]);

        translate([pcb_padding_x + pcb_screw_padding, pcb_padding_y + pcb_screw_padding, 0])
            for(i=[0:1])
                for(j=[0:1])
                    translate([i * pcb_screw_distance_x, j * pcb_screw_distance_y, thickness])
                        cylinder(pcb_screw_base_height, d=pcb_screw_base_d, $fn=20);
    }

    translate([pcb_padding_x + pcb_screw_padding, pcb_padding_y + pcb_screw_padding, 0]) {
        translate([pcb_pin0_x + button0_x, pcb_pin0_y + button0_y, 0])
            for(i=[0:3])
                for(j=[0:1])
                    translate([i * button_distance, j * button_distance, 0])
                        cylinder(thickness, d=12, $fn=20);

        for(i=[0:1])
            for(j=[0:1])
                translate([i * pcb_screw_distance_x, j * pcb_screw_distance_y, thickness])
                    cylinder(pcb_screw_base_height, d=pcb_screw_d - 0.2, $fn=20);
    }

    for (i=[0:1])
        for (j=[0:1])
            translate([i * (width - thickness),
                       screw_base_dim(screw_d) / 2 + j * (length - screw_base_dim(screw_d)),
                       height - screw_base_dim(screw_d) / 2])
                rotate([0, 90, 0])
                    cylinder(thickness, d=screw_d * 1.1, $fn=20);

    translate([thickness, (thickness - gap) / 2, thickness / 2])
        cube([width - 2 * thickness, thickness / 2 + gap, thickness / 2]);
    translate([thickness, length - thickness - gap / 2, thickness / 2])
        cube([width - 2 * thickness, thickness / 2 + gap, thickness / 2]);
}
