# SPDX-FileCopyrightText: 2022 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: GPL-2.0

cmake_minimum_required(VERSION 3.12)

include(CheckIPOSupported)
include(FetchContent)

FetchContent_Declare(avr_cmake
    GIT_REPOSITORY https://github.com/rafaelmartins/avr-cmake.git
    GIT_TAG 8741c4a95c9e03b0db56b8ca6cc7260259375e31
)
FetchContent_Declare(v_usb
    GIT_REPOSITORY https://github.com/obdev/v-usb.git
    GIT_TAG 7a28fdc685952412dad2b8842429127bc1cf9fa7
)
FetchContent_MakeAvailable(avr_cmake v_usb)

include(${avr_cmake_SOURCE_DIR}/cmake/avr.cmake)

project(b8 C ASM)

option(WITH_MCU "AVR Microcontroller.")
if(NOT WITH_MCU)
    set(WITH_MCU "attiny4313" CACHE STRING "AVR Microcontroller" FORCE)
endif()

add_executable(b8
    main.c
    bits.h
    usbconfig.h
    ${v_usb_SOURCE_DIR}/usbdrv/usbdrv.c
    ${v_usb_SOURCE_DIR}/usbdrv/usbdrv.h
    ${v_usb_SOURCE_DIR}/usbdrv/usbdrvasm.S
)

# the code for automatically generated usb serial number won't fit 2k rom
if(NOT WITH_MCU MATCHES "attiny2313*")
    target_sources(b8 PRIVATE
        serialnumber.c
        serialnumber.h
    )
endif()

target_include_directories(b8 PRIVATE
    .
    ${v_usb_SOURCE_DIR}
)

target_compile_options(b8 PRIVATE
    -Wall
    -Wextra
    -Werror
)

check_ipo_supported()
set_property(TARGET b8 PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE)

avr_target_set_device(b8 "${WITH_MCU}" 20000000UL)
avr_target_show_size(b8 "${WITH_MCU}")
avr_target_gen_map(b8)
avr_target_generate_hex(b8)
