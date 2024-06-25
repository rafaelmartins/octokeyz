# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

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
    OCTOKEYZ_VERSION=\"${OCTOKEYZ_VERSION}\"
)

set_property(TARGET octokeyz
    PROPERTY INTERPROCEDURAL_OPTIMIZATION TRUE
)

stm32f0_target_set_mcu(octokeyz stm32f042k6)
stm32f0_target_generate_map(octokeyz)
stm32f0_target_generate_bin(octokeyz)
stm32f0_target_generate_ihex(octokeyz)
stm32f0_target_generate_dfu(octokeyz)
stm32f0_target_show_size(octokeyz)
stm32f0_target_set_linker_script(octokeyz ${CMAKE_CURRENT_SOURCE_DIR}/STM32F042KxTx_FLASH.ld)
stm32f0_target_stlink_write(octokeyz)
stm32f0_add_stlink_targets()

install(FILES
    ${CMAKE_CURRENT_BINARY_DIR}/octokeyz.elf
    ${CMAKE_CURRENT_BINARY_DIR}/octokeyz.elf.map
    ${CMAKE_CURRENT_BINARY_DIR}/octokeyz.bin
    ${CMAKE_CURRENT_BINARY_DIR}/octokeyz.hex
    ${CMAKE_CURRENT_BINARY_DIR}/octokeyz.dfu
    TYPE BIN
)
