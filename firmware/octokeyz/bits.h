// SPDX-FileCopyrightText: 2022-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: GPL-2.0

#pragma once

#define P_LED (D, 6)

#define BIT_SET(reg, bit)   {(reg) |=  (1 << bit);}
#define BIT_CLEAR(reg, bit) {(reg) &= ~(1 << bit);}
#define BIT_FLIP(reg, bit)  {(reg) ^=  (1 << bit);}

#define _DDR_SET(port, bit)    BIT_SET(DDR ## port, bit)
#define _PORT_SET(port, bit)   BIT_SET(PORT ## port, bit)
#define _PORT_CLEAR(port, bit) BIT_CLEAR(PORT ## port, bit)
#define _PORT_FLIP(port, bit)  BIT_FLIP(PORT ## port, bit)

#define DDR_SET(p)    _DDR_SET p
#define PORT_SET(p)   _PORT_SET p
#define PORT_CLEAR(p) _PORT_CLEAR p
#define PORT_FLIP(p)  _PORT_FLIP p
