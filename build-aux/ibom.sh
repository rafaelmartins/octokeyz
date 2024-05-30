#!/bin/bash
# SPDX-FileCopyrightText: 2024 Rafael G. Martins <rafael@rafaelmartins.eng.br>
# SPDX-License-Identifier: BSD-3-Clause

set -e

MYDIR="$(realpath "$(dirname "${0}")")"
ROOTDIR="$(realpath "${MYDIR}/../")"
MYTMPDIR="$(mktemp -d)"

trap 'rm -rf -- "${MYTMPDIR}"' EXIT

if [[ x$CI = xtrue ]]; then
    sudo add-apt-repository -y ppa:kicad/kicad-8.0-releases
    sudo apt update
    sudo apt install -y --no-install-recommends kicad

    export PATH="${ROOTDIR}/InteractiveHtmlBom/InteractiveHtmlBom:${PATH}"
fi

export INTERACTIVE_HTML_BOM_NO_DISPLAY=1

function generate() {
    generate_interactive_bom.py \
        --no-browser \
        --dest-dir "${MYTMPDIR}" \
        --name-format "%f" \
        --include-tracks \
        --include-nets \
        --blacklist "H*" \
        "${ROOTDIR}/pcb/${1}/${1}.kicad_pcb"
}

generate octokeyz
generate octokeyz-mega

mv "${MYTMPDIR}"/*.html .
