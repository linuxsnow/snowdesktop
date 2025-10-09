#!/usr/bin/env bash

set -euo pipefail

OUTPUTDIR=mkosi.output
rm -f "$OUTPUTDIR/SHA256SUMS"

# find the first file in the output directory that matches the pattern "SNOW*_x86-64.manifest"
MANIFEST_FILE=$(find "$OUTPUTDIR" -maxdepth 1 -type f -name "SnowLinux*_x86-64.manifest" | head -n 1)
echo "Found manifest file: $MANIFEST_FILE"
SNOWVERSION=$(cat $MANIFEST_FILE | jq -r  '.config.version ')
echo "SNOWVERSION is: $SNOWVERSION"
rm -f "$OUTPUTDIR/v${SNOWVERSION}.SHA256SUMS"


# find all the files in the output directory that end in .SHA256SUMS and concatenate them into a single file named SHA256SUMS in the output directory
for file in "$OUTPUTDIR"/*.SHA256SUMS; do
    echo "Processing $file"
    cat "$file" >> "$OUTPUTDIR/v${SNOWVERSION}.SHA256SUMS"
done
