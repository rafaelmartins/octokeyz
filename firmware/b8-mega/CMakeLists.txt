# SPDX-FileCopyrightText: 2023-2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
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

project(b8-mega C ASM)

add_executable(b8-mega
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

target_link_libraries(b8-mega PRIVATE
    usbd-fs-stm32
)

target_compile_options(b8-mega PRIVATE
    -Wall
    -Wextra
    -Werror
)

set_property(TARGET b8-mega
    PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE
)

stm32f0_target_set_mcu(b8-mega stm32f042k6)
stm32f0_target_generate_map(b8-mega)
stm32f0_target_generate_ihex(b8-mega)
stm32f0_target_generate_dfu(b8-mega)
stm32f0_target_show_size(b8-mega)
stm32f0_target_set_linker_script(b8-mega ${CMAKE_CURRENT_SOURCE_DIR}/STM32F042K6Tx_FLASH.ld)
stm32f0_target_stlink_write(b8-mega)
stm32f0_add_stlink_targets()
