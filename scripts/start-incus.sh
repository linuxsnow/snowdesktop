#!/usr/bin/env bash
set -euo pipefail


# make the instance_name "snow" plus the variant
instance_name="snow-desktop"



incus start "$instance_name"
# this blocks until secure boot enrollment is complete
incus console --type=vga "$instance_name"
