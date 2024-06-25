# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

cmake_minimum_required(VERSION 3.22)

set(FETCHCONTENT_QUIET NO)
include(FetchContent)

FetchContent_Declare(cmake_stm32f0_cmsis
    GIT_REPOSITORY https://github.com/rafaelmartins/cmake-stm32f0-cmsis.git
    GIT_TAG f06fa52e39518ef5733fc333385535cbdf962399
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

find_package(Git)

if (EXISTS "${CMAKE_SOURCE_DIR}/version.cmake")
    include("${CMAKE_SOURCE_DIR}/version.cmake")
elseif(EXISTS "${CMAKE_SOURCE_DIR}/.git" AND Git_FOUND)
    execute_process(
        COMMAND ${GIT_EXECUTABLE} describe --abbrev=4 HEAD
        OUTPUT_VARIABLE _octokeyz_version
        ERROR_VARIABLE _octokeyz_version_err
        WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
        RESULT_VARIABLE _octokeyz_version_result
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )

    if(_octokeyz_version_err)
        message(FATAL_ERROR "Failed to find version from Git\n${_octokeyz_version_err}")
    endif()

    if(NOT _octokeyz_version_result EQUAL 0)
        message(FATAL_ERROR "Failed to find version from Git. Git process returned ${_octokeyz_version_result}")
    endif()

    string(REGEX REPLACE "^v" "" _octokeyz_version "${_octokeyz_version}")
    string(REGEX REPLACE "^([^-]*)-(.*)" "\\1.\\2" _octokeyz_version "${_octokeyz_version}")
    string(REGEX REPLACE "-g" "-" OCTOKEYZ_VERSION "${_octokeyz_version}")
else()
    message(FATAL_ERROR "Can't find version information!")
endif()

string(REGEX MATCHALL "[0-9]+" _octokeyz_version_list "${OCTOKEYZ_VERSION}")
list(LENGTH _octokeyz_version_list _octokeyz_version_list_len)
if(NOT _octokeyz_version_list_len GREATER_EQUAL 2)
    message(FATAL_ERROR "Invalid version: ${OCTOKEYZ_VERSION}")
endif()

list(GET _octokeyz_version_list 0 _octokeyz_version_major)
list(GET _octokeyz_version_list 1 _octokeyz_version_minor)
set(OCTOKEYZ_VERSION_CANONICAL "${_octokeyz_version_major}.${_octokeyz_version_minor}")

if(EXISTS "${CMAKE_SOURCE_DIR}/.git" AND _octokeyz_version_list_len GREATER_EQUAL 3)
    list(GET _octokeyz_version_list 2 _octokeyz_version_patch)
    set(OCTOKEYZ_VERSION_CANONICAL "${OCTOKEYZ_VERSION_CANONICAL}.${_octokeyz_version_patch}")
endif()

project(octokeyz
    VERSION "${OCTOKEYZ_VERSION_CANONICAL}"
    HOMEPAGE_URL "https://octokeyz.rgm.io/"
    LANGUAGES C ASM
)

add_subdirectory(firmware)

foreach(variant octokeyz octokeyz-mega)
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

set(CPACK_VERBATIM_VARIABLES YES)

set(CPACK_GENERATOR "TXZ" "ZIP")
set(CPACK_PACKAGE_VENDOR "Rafael G. Martins")
set(CPACK_PACKAGE_VERSION "${OCTOKEYZ_VERSION}")
set(CPACK_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}-${OCTOKEYZ_VERSION}")
set(CPACK_PACKAGE_CHECKSUM SHA512)

set(CPACK_INSTALL_SCRIPTS
    "${CMAKE_CURRENT_SOURCE_DIR}/cmake/cpack/install_license_and_readme.cmake"
    "${CMAKE_CURRENT_SOURCE_DIR}/cmake/cpack/install_version.cmake"
)

set(CPACK_SOURCE_GENERATOR "TXZ" "ZIP")
set(CPACK_SOURCE_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}-source-${OCTOKEYZ_VERSION}")
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
