// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#pragma once

#include <stdbool.h>
#include <stdint.h>

void idle_init(void);
void idle_set(uint8_t val);
uint8_t idle_get(void);
bool idle_request(void);
