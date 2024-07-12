# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

foreach(variant octokeyz octokeyz-mega)
    file(GLOB stl ${variant}/*.stl)
    install(FILES
        ${stl}
        DESTINATION share/3d-models/${variant}
    )
endforeach()
