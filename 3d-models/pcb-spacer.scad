/*
 * SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
 * SPDX-License-Identifier: CERN-OHL-S-2.0
 */

include <settings.scad>

difference() {
    cylinder(pcbs_spacer_height, d=pcb_screw_base_d, $fn=20);
    cylinder(pcbs_spacer_height, d=pcb_screw_d + 1, $fn=20);
}
