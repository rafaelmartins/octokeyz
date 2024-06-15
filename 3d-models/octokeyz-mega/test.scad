// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: CERN-OHL-S-2.0

// this is just an assembly of the STLs, so we can validate the sizes and positioning.
// no need to print anything.

include <settings.scad>

// pcb placeholder
translate([thickness + pcb_gap,
           thickness + pcb_gap,
           display_height - thickness - display_pcb_base_height - pcb_thickness]) {
    cube([switches_pcb_length, switches_pcb_width, pcb_thickness]);
    translate([switches_pcb_length - display_pcb_length, switches_pcb_width, 0])
        cube([display_pcb_length, display_pcb_width, pcb_thickness]);
}

%rotate([180, 0, 0])
    translate([0, -base_width, 0])
        import("base.stl");

%rotate([0, 0, 180])
    translate([-switches_length, -base_width + thickness, 0])
        import("enclosure-back.stl");

rotate([180, 0, 0])
    translate([0, -switches_width, -switches_height])
        import("enclosure-front-switches.stl");

rotate([180, 0, 0])
    translate([switches_length - display_length, -base_width, -display_height])
        import("enclosure-front-display.stl");