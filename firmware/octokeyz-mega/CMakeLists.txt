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
    GIT_TAG 481ec4e92ece803a0c6f3704059bcba6318f3676
)
FetchContent_MakeAvailable(
    cmake_stm32f0_cmsis
    usbd_fs_stm32
)

include(${cmake_stm32f0_cmsis_SOURCE_DIR}/cmake/stm32f0-cmsis.cmake)

project(octokeyz-mega C ASM)

add_executable(octokeyz-mega
    descriptors.c
    display.c
    display.h
    display-font.h
    led.c
    led.h
    main.c
    usb-watchdog.c
    usb-watchdog.h
)

target_compile_definitions(usbd-fs-stm32 INTERFACE
    USBD_EP1_IN_SIZE=64
    USBD_EP1_OUT_SIZE=64
    USBD_EP1_TYPE=INTERRUPT
)

target_link_libraries(octokeyz-mega PRIVATE
    usbd-fs-stm32
)

target_compile_options(octokeyz-mega PRIVATE
    -Wall
    -Wextra
    -Werror
)

set_property(TARGET octokeyz-mega
    PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE
)

stm32f0_target_set_mcu(octokeyz-mega stm32f042k6)
stm32f0_target_generate_map(octokeyz-mega)
stm32f0_target_generate_ihex(octokeyz-mega)
stm32f0_target_generate_dfu(octokeyz-mega)
stm32f0_target_show_size(octokeyz-mega)
stm32f0_target_set_linker_script(octokeyz-mega ${CMAKE_CURRENT_SOURCE_DIR}/STM32F042K6Tx_FLASH.ld)
stm32f0_target_stlink_write(octokeyz-mega)
stm32f0_add_stlink_targets()
