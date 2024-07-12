# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

find_package(Git)

if (EXISTS "${CMAKE_SOURCE_DIR}/version.cmake")
    include("${CMAKE_SOURCE_DIR}/version.cmake")
elseif(EXISTS "${CMAKE_SOURCE_DIR}/.git" AND Git_FOUND)
    execute_process(
        COMMAND ${GIT_EXECUTABLE} describe --abbrev=4 HEAD
        OUTPUT_VARIABLE _gvg_version
        ERROR_VARIABLE _gvg_version_err
        WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
        RESULT_VARIABLE _gvg_version_result
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )

    if(_gvg_version_err)
        message(FATAL_ERROR "Failed to find version from Git\n${_gvg_version_err}")
    endif()

    if(NOT _gvg_version_result EQUAL 0)
        message(FATAL_ERROR "Failed to find version from Git. Git process returned ${_gvg_version_result}")
    endif()

    string(REGEX REPLACE "^v" "" _gvg_version "${_gvg_version}")
    string(REGEX REPLACE "^([^-]*)-(.*)" "\\1.\\2" _gvg_version "${_gvg_version}")
    string(REGEX REPLACE "-g" "-" PACKAGE_VERSION "${_gvg_version}")
else()
    message(FATAL_ERROR "Can't find version information!")
endif()

string(REGEX MATCHALL "[0-9]+" _gvg_version_list "${PACKAGE_VERSION}")
list(LENGTH _gvg_version_list _gvg_version_list_len)
if(NOT _gvg_version_list_len GREATER_EQUAL 2)
    message(FATAL_ERROR "Invalid version: ${PACKAGE_VERSION}")
endif()

list(GET _gvg_version_list 0 _gvg_version_major)
list(GET _gvg_version_list 1 _gvg_version_minor)
set(PACKAGE_VERSION_CANONICAL "${_gvg_version_major}.${_gvg_version_minor}")

if(EXISTS "${CMAKE_SOURCE_DIR}/.git" AND _gvg_version_list_len GREATER_EQUAL 3)
    list(GET _gvg_version_list 2 _gvg_version_patch)
    set(PACKAGE_VERSION_CANONICAL "${PACKAGE_VERSION_CANONICAL}.${_gvg_version_patch}")
endif()
