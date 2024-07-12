# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

cmake_minimum_required(VERSION 3.22)

set(FETCHCONTENT_QUIET NO)
include(FetchContent)

FetchContent_Declare(cmake_cmsis_stm32
    GIT_REPOSITORY https://github.com/rafaelmartins/cmake-cmsis-stm32.git
    GIT_TAG fe2a4affff8c36d431faf5f9890bac811a1578da
)
FetchContent_Declare(usbd_fs_stm32
    GIT_REPOSITORY https://github.com/rafaelmartins/usbd-fs-stm32.git
    GIT_TAG 94db86e1b47dd830fab0ea736d006a0b3723189d
)
FetchContent_MakeAvailable(
    cmake_cmsis_stm32
    usbd_fs_stm32
)

include(cmake/git-version-gen.cmake)

project(octokeyz
    VERSION "${PACKAGE_VERSION_CANONICAL}"
    HOMEPAGE_URL "https://octokeyz.rgm.io/"
    LANGUAGES C ASM
)

find_program(GENERATE_INTERACTIVE_BOM
    generate_interactive_bom.py
    HINTS ${CMAKE_SOURCE_DIR}/../InteractiveHtmlBom/InteractiveHtmlBom
)

if(GENERATE_INTERACTIVE_BOM)
    add_custom_target(ibom ALL)
endif()

foreach(variant octokeyz octokeyz-mega)
    if(GENERATE_INTERACTIVE_BOM)
        add_custom_command(
            OUTPUT ibom/${variant}.html
            DEPENDS "${CMAKE_SOURCE_DIR}/pcb/${variant}/${variant}.kicad_pcb"
            COMMAND
                generate_interactive_bom.py
                    --no-browser
                    --dest-dir "${CMAKE_BINARY_DIR}/ibom"
                    --name-format "%f"
                    --include-tracks
                    --include-nets
                    --blacklist "H*"
                    "${CMAKE_SOURCE_DIR}/pcb/${variant}/${variant}.kicad_pcb"
        )
        add_custom_target(ibom-${variant}
            DEPENDS ibom/${variant}.html
        )
        install(FILES
            ${CMAKE_BINARY_DIR}/ibom/${variant}.html
            DESTINATION share/pcb/${variant}
        )
        add_dependencies(ibom
            ibom-${variant}
        )
    endif()

    file(GLOB gerber pcb/${variant}/gerber/${variant}_rev*.zip)
    install(FILES
        pcb/${variant}/${variant}.pdf
        ${gerber}
        DESTINATION share/pcb/${variant}
    )

    file(GLOB stl 3d-models/${variant}/*.stl)
    install(FILES
        ${stl}
        DESTINATION share/3d-models/${variant}
    )
endforeach()

install(DIRECTORY
    share/images
    share/udev
    DESTINATION share
)

add_subdirectory(firmware)

set(CPACK_VERBATIM_VARIABLES YES)

set(CPACK_GENERATOR "TXZ" "ZIP")
set(CPACK_PACKAGE_VENDOR "Rafael G. Martins")
set(CPACK_PACKAGE_VERSION "${PACKAGE_VERSION}")
set(CPACK_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}-${PACKAGE_VERSION}")
set(CPACK_PACKAGE_CHECKSUM SHA512)

set(CPACK_INSTALL_SCRIPTS
    "${CMAKE_SOURCE_DIR}/cmake/cpack/install_license_and_readme.cmake"
    "${CMAKE_SOURCE_DIR}/cmake/cpack/install_version.cmake"
)

set(CPACK_SOURCE_GENERATOR "TXZ" "ZIP")
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
