# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

cmake_minimum_required(VERSION 3.22)

include(FetchContent)

FetchContent_Declare(cmake_cmsis_stm32
    GIT_REPOSITORY https://github.com/rafaelmartins/cmake-cmsis-stm32.git
    GIT_TAG fe2a4affff8c36d431faf5f9890bac811a1578da
)
FetchContent_MakeAvailable(
    cmake_cmsis_stm32
)

include(cmake/git-version-gen.cmake)

project(octokeyz
    VERSION "${PACKAGE_VERSION_CANONICAL}"
    HOMEPAGE_URL "https://octokeyz.rgm.io/"
    LANGUAGES C ASM
)

add_subdirectory(3d-models)
add_subdirectory(firmware)
add_subdirectory(pcb)
add_subdirectory(share)

configure_file(
    cmake/cpack/files/README.txt.in
    cmake/cpack/files/README.txt
)

install(FILES
    cmake/cpack/files/LICENSE.txt
    "${CMAKE_BINARY_DIR}/cmake/cpack/files/README.txt"
    DESTINATION .
)

set(CPACK_VERBATIM_VARIABLES YES)

set(CPACK_GENERATOR TXZ ZIP)
set(CPACK_PACKAGE_VENDOR "Rafael G. Martins")
set(CPACK_PACKAGE_VERSION "${PACKAGE_VERSION}")
set(CPACK_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}-${PACKAGE_VERSION}")
set(CPACK_PACKAGE_CHECKSUM SHA512)

set(CPACK_INSTALL_SCRIPTS
    "${CMAKE_SOURCE_DIR}/cmake/cpack/install_version.cmake"
)

set(CPACK_SOURCE_GENERATOR TXZ ZIP)
set(CPACK_SOURCE_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}-source-${PACKAGE_VERSION}")
set(CPACK_SOURCE_IGNORE_FILES
    \\.git
    \\.gitignore
    \\.github/
    "build.*/"
    "/.*-backups/"
    "/gerber/"
    "-cache$"
    "-cache.lib$"
    "\\.(000|bak|bck|kicad_pcb-bak|lck|tmp|dsn|ses|kicad_sch-bak|kicad_prl)$"
    ".*~$"
)

include(CPack)
