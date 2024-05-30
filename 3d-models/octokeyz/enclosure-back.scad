// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: CERN-OHL-S-2.0

include <../lib/screw-base.scad>
include <settings.scad>

length_ = length - 2 * (thickness + gap);
height_ = height - thickness - gap;

difference() {
    union() {
        cube([width, length_, thickness]);
        cube([thickness, length_, height_]);
        translate([width - thickness, 0, 0])
            cube([thickness, length_, height_]);

        rotate([-90, -90, 0])
            screw_base_add(screw_d, screw_h);
        translate([0, length_, 0])
            rotate([90, 0, 0])
                screw_base_add(screw_d, screw_h);
        translate([width, 0, 0])
            rotate([-90, 180, 0])
                screw_base_add(screw_d, screw_h);
        translate([width, length_, 0])
            rotate([90, -90, 0])
                screw_base_add(screw_d, screw_h);

        translate([thickness / 2, 4, height_])
            cube([thickness / 2, length_ - 8, thickness / 2]);
        translate([width - thickness, 4, height_])
            cube([thickness / 2, length_ - 8, thickness / 2]);

        translate([width - usb_distance_x - usb_width + gap, -(thickness + gap), 0])
            cube([usb_width - 2 * gap, thickness + gap, height - usb_distance_z - usb_height]);
    }

    rotate([-90, -90, 0])
        screw_base_sub(screw_d, screw_h);
    translate([0, length_, 0])
        rotate([90, 0, 0])
            screw_base_sub(screw_d, screw_h);
    translate([width, 0, 0])
        rotate([-90, 180, 0])
            screw_base_sub(screw_d, screw_h);
    translate([width, length_, 0])
        rotate([90, -90, 0])
            screw_base_sub(screw_d, screw_h);
}
