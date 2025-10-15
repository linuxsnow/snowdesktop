#!/bin/bash

set -euo pipefail

download_flatpaks() {

    # Only check out en. We don't really support other languages on the live image at this time.
    flatpak config --set languages en

    # read all the flatpaks from /usr/share/flatpak/install-list and install them
    if [ -f /usr/share/flatpak/install-list ]; then
        mapfile -t flatpaks < /usr/share/flatpak/install-list
        for flatpak in "${flatpaks[@]}"; do
            echo "Installing flatpak: $flatpak"
            flatpak install --or-update --noninteractive --assumeyes flathub "$flatpak"
        done
    fi


}

# bail out if root=tmpfs on the command line, as this is a live image
if grep -q 'root=tmpfs' /proc/cmdline; then
    echo "Skipping flatpak installation as this is a live image (root=tmpfs)"
    exit 0
fi

# check for a marker file in /etc/snow to see if we should skip this
if [ -f /etc/snow/flatpaks-installed ]; then
    echo "Skipping flatpak installation as /etc/snow/flatpaks-installed exists"
    exit 0
fi
download_flatpaks
mkdir -p /etc/snow
echo "Flatpak installation complete" > /etc/snow/flatpak-installed

