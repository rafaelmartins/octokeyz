// SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#pragma once

void usb_watchdog_init(void);
void usb_watchdog_reset(void);
void usb_watchdog_task(void);

// callbacks
void usb_watchdog_cb(void);
