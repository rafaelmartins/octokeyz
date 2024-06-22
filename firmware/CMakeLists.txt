# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

cmake_minimum_required(VERSION 3.17)

set(FETCHCONTENT_QUIET NO)
include(FetchContent)

FetchContent_Declare(cmake_stm32f0_cmsis
    GIT_REPOSITORY https://github.com/rafaelmartins/cmake-stm32f0-cmsis.git
    GIT_TAG 27f78dce1bee70f12ce3aee1eade57301b993938
)
FetchContent_Declare(usbd_fs_stm32
    GIT_REPOSITORY https://github.com/rafaelmartins/usbd-fs-stm32.git
    GIT_TAG 39c612119d527b7db4efcd5964440e6a906f20c5
)
FetchContent_MakeAvailable(
    cmake_stm32f0_cmsis
    usbd_fs_stm32
)

include(${cmake_stm32f0_cmsis_SOURCE_DIR}/cmake/stm32f0-cmsis.cmake)

project(octokeyz C ASM)

add_executable(octokeyz
    bootloader.c
    bootloader.h
    descriptors.c
    display.c
    display.h
    display-font.h
    led.c
    led.h
    main.c
)

target_compile_definitions(usbd-fs-stm32 INTERFACE
    USBD_EP1_IN_SIZE=64
    USBD_EP1_OUT_SIZE=64
    USBD_EP1_TYPE=INTERRUPT
)

target_link_libraries(octokeyz PRIVATE
    usbd-fs-stm32
)

target_compile_options(octokeyz PRIVATE
    -Wall
    -Wextra
    -Werror
)

set_property(TARGET octokeyz
    PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE
)

stm32f0_target_set_mcu(octokeyz stm32f042k6)
stm32f0_target_generate_map(octokeyz)
stm32f0_target_generate_ihex(octokeyz)
stm32f0_target_generate_dfu(octokeyz)
stm32f0_target_show_size(octokeyz)
stm32f0_target_set_linker_script(octokeyz ${CMAKE_CURRENT_SOURCE_DIR}/STM32F042K6Tx_FLASH.ld)
stm32f0_target_stlink_write(octokeyz)
stm32f0_add_stlink_targets()
