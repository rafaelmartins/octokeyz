# SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: GPL-2.0

cmake_minimum_required(VERSION 3.12)

include(CheckIPOSupported)
include(./cmake/avr.cmake)

project(b8 C ASM)

add_executable(b8
    main.c
    bits.h
    usbconfig.h
    usbdrv/usbdrv.c
    usbdrv/usbdrv.h
    usbdrv/usbdrvasm.S
)

target_include_directories(b8 PRIVATE
    .
)

target_compile_options(b8 PRIVATE
    -Wall
    -Wextra
    -Werror
)

check_ipo_supported()
set_property(TARGET b8 PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE)

avr_target_set_device(b8 attiny2313a 20000000UL)
avr_target_show_size(b8 attiny2313a)
avr_target_gen_map(b8)
avr_target_generate_hex(b8)
