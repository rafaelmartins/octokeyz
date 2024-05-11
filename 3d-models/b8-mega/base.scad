/*
 * SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

// should print upside down

include <settings.scad>

display_width_ = display_width - 2 * (thickness + gap);
switches_width_ = switches_width - 2 * (thickness + gap);


difference() {
    union() {
        cube([base_length, base_width, thickness]);

        translate([0, 0, thickness])
            rotate([90, 0, 90])
            linear_extrude(height=base_length)
                polygon([[0, 0], [base_width, 0], [0, base_height - thickness]]);
    }

    cube([base_length - display_length, display_width_ + 2 * gap + pcb_gap, base_height]);

    translate([base_length - display_length, thickness + gap, 0]) {
        translate([base_screw_padding, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
        translate([display_length - base_screw_padding, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
    }

    translate([0, display_width_ + thickness + gap, 0]) {
        translate([base_screw_padding, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
        translate([switches_length / 2, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
        translate([switches_length - base_screw_padding, base_screw_padding, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
        translate([base_screw_padding, switches_width_ - base_screw_padding_bottom, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
        translate([switches_length / 2, switches_width_ - base_screw_padding_bottom, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
        translate([switches_length - base_screw_padding, switches_width_ - base_screw_padding_bottom, 0])
            cylinder(h=base_screw_h, d=base_screw_d, $fn=20);
    }
}
