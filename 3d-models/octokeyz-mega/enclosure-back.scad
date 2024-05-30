// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: CERN-OHL-S-2.0

include <../lib/screw-base.scad>
include <settings.scad>

display_width_ = display_width - 2 * (thickness + gap);
display_height_ = display_height - thickness - gap;

switches_width_ = switches_width - 2 * (thickness + gap);
switches_height_ = switches_height - thickness - gap;


difference() {
    union() {
        cube([display_length, display_width_, thickness]);
        cube([thickness, display_width_, display_height_]);
        translate([display_length - thickness, 0, 0])
            cube([thickness, display_width_, display_height_]);

        translate([thickness / 2, 4, display_height_])
            cube([thickness / 2, display_width_ - 8, thickness / 2]);
        translate([display_length - thickness, 4, display_height_])
            cube([thickness / 2, display_width_ - 8, thickness / 2]);

        translate([display_length - usb_distance_x - usb_width + gap, -thickness - gap, 0])
            cube([usb_width - 2 * gap, thickness + gap, display_height - usb_distance_z - usb_height]);

        rotate([-90, -90, 0])
            screw_base_add(screw_d, screw_h);
        translate([display_length, 0, 0])
            rotate([-90, 180, 0])
                screw_base_add(screw_d, screw_h);

        translate([0, display_width_, 0]) {
            cube([display_length, 2 * gap + pcb_gap, thickness]);
            cube([thickness, 2 * gap + pcb_gap, switches_height_]);

            translate([0, 2 * gap + pcb_gap, 0]) {
                cube([switches_length, switches_width_, thickness]);
                cube([thickness, switches_width_, switches_height_]);
                translate([switches_length - thickness, 0, 0])
                    cube([thickness, switches_width_, switches_height_]);

                translate([thickness / 2, 4, switches_height_])
                    cube([thickness / 2, switches_width_ - 8, thickness / 2]);
                translate([switches_length - thickness, 4, switches_height_])
                    cube([thickness / 2, switches_width_ - 8, thickness / 2]);

                translate([display_length + back_screw3_padding, 0, 0])
                    rotate([-90, -90, 0])
                        screw_base_add(screw_d, screw_h);
                translate([switches_length, 0, 0])
                    rotate([-90, 180, 0])
                        screw_base_add(screw_d, screw_h);

                translate([0, switches_width_, 0])
                    rotate([90, 0, 0])
                        screw_base_add(screw_d, screw_h);
                translate([switches_length, switches_width_, 0])
                    rotate([90, -90, 0])
                        screw_base_add(screw_d, screw_h);
            }
        }
    }

    rotate([-90, -90, 0])
        screw_base_sub(screw_d, screw_h);
    translate([display_length, 0, 0])
        rotate([-90, 180, 0])
            screw_base_sub(screw_d, screw_h);

    translate([0, display_width_ + 2 * gap + pcb_gap, 0]) {
        translate([display_length + back_screw3_padding, -0.1, 0]) // workaround
            rotate([-90, -90, 0])
                screw_base_sub(screw_d, screw_h);
        translate([switches_length, -0.1, 0]) // workaround
            rotate([-90, 180, 0])
                screw_base_sub(screw_d, screw_h);
        translate([0, switches_width_, 0])
            rotate([90, 0, 0])
                screw_base_sub(screw_d, screw_h);
        translate([switches_length, switches_width_, 0])
            rotate([90, -90, 0])
                screw_base_sub(screw_d, screw_h);
    }

    translate([base_screw_padding, base_screw_padding, 0])
        cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
    translate([display_length - base_screw_padding, base_screw_padding, 0])
        cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);

    translate([0, display_width_, 0]) {
        translate([base_screw_padding, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
        translate([switches_length / 2, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
        translate([switches_length - base_screw_padding, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
        translate([base_screw_padding, switches_width_ - base_screw_padding_bottom, 0])
            cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
        translate([switches_length / 2, switches_width_ - base_screw_padding_bottom, 0])
            cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
        translate([switches_length - base_screw_padding, switches_width_ - base_screw_padding_bottom, 0])
            cylinder(h=base_screw_h, d=base_screw_d * 1.1, $fn=20);
    }
}
