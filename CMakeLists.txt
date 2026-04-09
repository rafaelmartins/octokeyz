# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

cmake_minimum_required(VERSION 3.25)

include(FetchContent)

FetchContent_Declare(cmake_cmsis_stm32
    GIT_REPOSITORY https://github.com/rafaelmartins/cmake-cmsis-stm32.git
    GIT_TAG c199670047e5adb373d8a033ece58c34f173ad4b
)
FetchContent_Declare(cmake_stm32_firmware_release
    GIT_REPOSITORY https://github.com/rafaelmartins/cmake-stm32-firmware-release.git
    GIT_TAG 983970487fad111d48e7b7dca9e0e840e055705a
)
FetchContent_MakeAvailable(
    cmake_cmsis_stm32
    cmake_stm32_firmware_release
)

project(octokeyz
    VERSION "${PACKAGE_VERSION_CANONICAL}"
    HOMEPAGE_URL "https://rafaelmartins.com/p/octokeyz/"
    LANGUAGES C ASM
)

add_subdirectory(firmware)

cpack_configure()
