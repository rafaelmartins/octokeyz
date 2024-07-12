# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

option(GENERATE_IBOM "Enable Interactive Bill of Materials. Requires kicad and InteractiveHtmlBom")
if(GENERATE_IBOM)
    find_program(GIB
        generate_interactive_bom.py
        HINTS "${CMAKE_SOURCE_DIR}/../InteractiveHtmlBom/InteractiveHtmlBom"
        REQUIRED
    )
    add_custom_target(ibom ALL)
endif()

foreach(variant octokeyz octokeyz-mega)
    if(GENERATE_IBOM AND GIB)
        add_custom_command(
            OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/${variant}.html
            DEPENDS "${CMAKE_SOURCE_DIR}/pcb/${variant}/${variant}.kicad_pcb"
            COMMAND
                ${GIB}
                    --no-browser
                    --dest-dir "${CMAKE_CURRENT_BINARY_DIR}"
                    --name-format "%f"
                    --include-tracks
                    --include-nets
                    --blacklist "H*"
                    "${CMAKE_CURRENT_SOURCE_DIR}/${variant}/${variant}.kicad_pcb"
        )
        add_custom_target(ibom-${variant}
            DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/${variant}.html
        )
        install(FILES
            ${CMAKE_CURRENT_BINARY_DIR}/${variant}.html
            DESTINATION share/pcb/${variant}
        )
        add_dependencies(ibom
            ibom-${variant}
        )
    endif()

    file(GLOB gerber ${variant}/gerber/${variant}_rev*.zip)
    install(FILES
        ${variant}/${variant}.pdf
        ${gerber}
        DESTINATION share/pcb/${variant}
    )

    file(READ ${variant}/${variant}.kicad_pcb kicad_pcb)
    if(NOT kicad_pcb)
        message(FATAL_ERROR "Failed to read ${variant}.kicad_pcb")
    endif()

    string(REGEX REPLACE ".*rev \"([0-9.]*)\".*" "\\1" rev "${kicad_pcb}")
    set(${variant}_PCB_REVISION "${rev}" PARENT_SCOPE)
endforeach()
