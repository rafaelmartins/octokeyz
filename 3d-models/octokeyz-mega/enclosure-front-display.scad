// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: CERN-OHL-S-2.0

include <../lib/screw-base.scad>
include <../lib/ssd1306.scad>
include <settings.scad>


difference() {
    union() {
        cube([display_length, thickness, display_height]);
        cube([display_length, display_width, thickness]);

        translate([oled0_x, oled0_y, 0])
            ssd1306_add(thickness);

        translate([0, display_width - thickness, thickness]) {
            cube([display_length, thickness, display_height_in]);

            translate([front_screw_padding, -screw_base_dim(screw_d) / 2,
                       screw_base_height(display_height_in - 1)])
                rotate([180, 0, 90])
                    screw_base_add(screw_d, display_height_in - 1);
            translate([display_length - front_screw_padding - screw_base_dim(screw_d),
                       -screw_base_dim(screw_d) / 2, screw_base_height(display_height_in - 1)])
                rotate([180, 0, 90])
                    screw_base_add(screw_d, display_height_in - 1);
        }

        translate([display_pcb0_x, display_pcb0_y, 0]) {
            translate([pcb_screw_padding, pcb_screw_padding, thickness]) {
                cylinder(h=display_pcb_base_height, d=pcb_base_d, $fn=20);
                translate([display_pcb_base_distance, 0, 0])
                    cylinder(h=display_pcb_base_height, d=pcb_base_d, $fn=20);
            }
        }
    }

    translate([oled0_x, oled0_y, 0])
        ssd1306_sub(thickness);

    translate([0, display_width - thickness, thickness]) {
        translate([front_screw_padding, -screw_base_dim(screw_d) / 2,
                   screw_base_height(display_height_in - 1)])
            rotate([180, 0, 90])
                screw_base_sub(screw_d, display_height_in - 1);
        translate([display_length - front_screw_padding - screw_base_dim(screw_d),
                   -screw_base_dim(screw_d) / 2, screw_base_height(display_height_in - 1)])
            rotate([180, 0, 90])
                screw_base_sub(screw_d, display_height_in - 1);
    }

    translate([display_pcb0_x, display_pcb0_y, 0]) {
        translate([pcb_screw_padding, pcb_screw_padding, thickness]) {
            cylinder(h=display_pcb_base_height, d=pcb_base_screw_d, $fn=20);
            translate([display_pcb_base_distance, 0, 0])
                cylinder(h=display_pcb_base_height, d=pcb_base_screw_d, $fn=20);
        }

        translate([led0_x, led0_y, 0])
            cylinder(h=thickness, d=led_d, $fn=20);
    }

    translate([usb_distance_x, 0, usb_distance_z])
        cube([usb_width, thickness, display_height - usb_distance_z]);

    translate([(thickness - gap) / 2, thickness, thickness / 2])
        cube([thickness / 2 + gap, display_width - 2 * thickness, thickness / 2]);
    translate([display_length - thickness - gap / 2, thickness, thickness / 2])
        cube([thickness / 2 + gap, display_width - 2 * thickness, thickness / 2]);

    translate([screw_base_dim(screw_d) / 2, 0, display_height - screw_base_dim(screw_d) / 2])
        rotate([-90, 0, 0])
            cylinder(h=thickness, d=screw_d * 1.1, $fn=20);
    translate([display_length - screw_base_dim(screw_d) / 2, 0,
               display_height - screw_base_dim(screw_d) / 2])
        rotate([-90, 0, 0])
            cylinder(h=thickness, d=screw_d * 1.1, $fn=20);
}
