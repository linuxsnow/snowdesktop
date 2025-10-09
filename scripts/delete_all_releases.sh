#!/bin/bash


gh release list \
  | awk -F '\t' '{print $3}' \
  | while read -r line; do
  gh release delete --cleanup-tag -y "$line"
done
