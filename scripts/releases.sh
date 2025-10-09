#!/usr/bin/env bash
set -euo pipefail

echo "Cleaning up releases..."

# get a list of all releases on the remote server
releases=$(ls /mnt/fast/snow/v*.SHA256SUMS || true)

# keep only the last 3 releases
releases_to_delete=$(echo "$releases" | sort -V | head -n -5)
if [[ -n "$releases_to_delete" ]]; then
    echo "Deleting old releases:"
    echo "$releases_to_delete"
    # for each file in releases_to_delete, cat the file to get the list of files to delete
    # the files are in the format "<sha256sum>  *<filename>"
    # we want to extract the filename and delete it from the remote server
    for release in $releases_to_delete; do
        echo "Processing $release"
        files_to_delete=$(cat $release | awk '{print $2}' | sed 's/^\*//')
        for file in $files_to_delete; do
            echo "Deleting $file ..."
            rm -f /mnt/fast/snow/$file
        done
        echo "Deleting release file $release ..."
        rm -f $release
    done
fi

keep_releases=$(ls /mnt/fast/snow/v*.SHA256SUMS || true)
echo "Keeping these releases:"
echo "$keep_releases"

# move previous SHA256SUMS files to .old
mv /mnt/fast/snow/SHA256SUMS /mnt/fast/snow/SHA256SUMS.old || true

# concatenate all the v*.SHA256SUMS files into a single SHA256SUMS file
cat /mnt/fast/snow/v*.SHA256SUMS > /mnt/fast/snow/SHA256SUMS
