/*
 * SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

include <../lib/screw-base.scad>
include <../lib/ssd1306.scad>
include <settings.scad>


difference() {
    union() {
        cube([switches_length - display_length, thickness, switches_height]);
        cube([switches_length - display_length, pcb_padding_y - thickness, thickness]);
        translate([0, pcb_padding_y - thickness, 0])
            cube([switches_length, switches_width - pcb_padding_y + thickness, thickness]);
        translate([0, switches_width - thickness, 0])
            cube([switches_length, thickness, switches_height]);

        translate([0, pcb_padding_y - screw_base_dim(screw_d), 0]) {
            translate([switches_length - display_length + front_screw_padding, 0, 0])
                cube([screw_base_dim(screw_d), screw_base_dim(screw_d), thickness]);
            translate([switches_length - front_screw_padding - screw_base_dim(screw_d), 0, 0])
                cube([screw_base_dim(screw_d), screw_base_dim(screw_d), thickness]);
        }

        translate([pcb_padding_x, pcb_padding_y, 0]) {
            translate([pcb_screw_padding, pcb_screw_padding, thickness])
                for (i=[0:2])
                    for (j=[0:2])
                        if (!(i == 1 && j == 1))
                            translate([i * switches_pcb_base_distance_x,
                                       j * switches_pcb_base_distance_y, 0])
                                cylinder(h=switches_pcb_base_height, d=pcb_base_d, $fn=20);
        }
    }

    translate([0, pcb_padding_y - screw_base_dim(screw_d), 0]) {
        translate([switches_length - display_length + front_screw_padding + screw_base_dim(screw_d) / 2,
                   screw_base_dim(screw_d) / 2, 0])
            cylinder(h=thickness, d=screw_d * 1.1, $fn=20);
        translate([switches_length - front_screw_padding - screw_base_dim(screw_d) / 2,
                   screw_base_dim(screw_d) / 2, 0])
            cylinder(h=thickness, d=screw_d * 1.1, $fn=20);
    }

    translate([pcb_padding_x, pcb_padding_y, 0]) {
        translate([pcb_screw_padding, pcb_screw_padding, thickness])
            for (i=[0:2])
                for (j=[0:2])
                    if (!(i == 1 && j == 1))
                        translate([i * switches_pcb_base_distance_x,
                                   j * switches_pcb_base_distance_y, 0])
                            cylinder(h=switches_pcb_base_height, d=pcb_base_screw_d, $fn=20);

         translate([key0_x, key0_y, 0])
             for (i=[0:3])
                 for (j=[0:1])
                     translate([i * key_distance, j * key_distance, 0])
                        cube([key_dim, key_dim, thickness]);
    }

    translate([(thickness - gap) / 2, thickness, thickness / 2])
        cube([thickness / 2 + gap, switches_width - 2 * thickness, thickness / 2]);
    translate([switches_length - thickness - gap / 2, pcb_padding_y - thickness, thickness / 2])
        cube([thickness / 2 + gap, switches_width - 2 * thickness + pcb_padding_y - thickness,
              thickness / 2]);

    translate([screw_base_dim(screw_d) / 2, 0, switches_height - screw_base_dim(screw_d) / 2])
        rotate([-90, 0, 0])
            cylinder(h=thickness, d=screw_d * 1.1, $fn=20);
    translate([(switches_length - display_length) - screw_base_dim(screw_d) / 2, 0,
               switches_height - screw_base_dim(screw_d) / 2])
        rotate([-90, 0, 0])
            cylinder(h=thickness, d=screw_d * 1.1, $fn=20);
    translate([0, switches_width - thickness - 0.1, 0]) {  // workaround
        translate([screw_base_dim(screw_d) / 2, 0, switches_height - screw_base_dim(screw_d) / 2])
            rotate([-90, 0, 0])
                cylinder(h=thickness + 0.2, d=screw_d * 1.1, $fn=20);  // workaround
        translate([switches_length - screw_base_dim(screw_d) / 2, 0,
                   switches_height - screw_base_dim(screw_d) / 2])
            rotate([-90, 0, 0])
                cylinder(h=thickness + 0.2, d=screw_d * 1.1, $fn=20);  // workaround
    }
}
