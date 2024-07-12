# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

include(FetchContent)

FetchContent_Declare(usbd_fs_stm32
    GIT_REPOSITORY https://github.com/rafaelmartins/usbd-fs-stm32.git
    GIT_TAG 94db86e1b47dd830fab0ea736d006a0b3723189d
)
FetchContent_MakeAvailable(
    usbd_fs_stm32
)

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

target_compile_definitions(octokeyz PRIVATE
    PACKAGE_VERSION=\"${PACKAGE_VERSION}\"
)

set_property(TARGET octokeyz
    PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE
)

cmsis_stm32_target(octokeyz
    DEVICE STM32F042x6
    VERSION 20240709193138
    LINKER_SCRIPT ${CMAKE_CURRENT_SOURCE_DIR}/STM32F042KxTx_FLASH.ld
    ADDITIONAL_OUTPUTS BIN MAP IHEX DFU
    SHOW_SIZE ON
    STLINK ON
    INSTALL ON
)
