# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

if(NOT CPACK_SOURCE_INSTALLED_DIRECTORIES)
    file(COPY_FILE
        "${CMAKE_CURRENT_LIST_DIR}/LICENSE"
        "${CMAKE_INSTALL_PREFIX}/LICENSE.txt"
    )

    file(COPY_FILE
        "${CMAKE_CURRENT_LIST_DIR}/../../README.md"
        "${CMAKE_INSTALL_PREFIX}/README.md"
    )
endif()
