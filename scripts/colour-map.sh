#!/bin/sh

url="https://api.brickowl.com/v1/catalog/color_list?key=${BRICKOWL_API_KEY}"
buildMap='[.[] | select(.ldraw_ids[0] != null)] | map({ (.ldraw_ids[0]): .hex}) | add'

json=$(curl "$url" --silent | jq --tab --sort-keys "$buildMap")


printf "package lego\n\nvar hexColours = \`%s\`\n" "$json" > ./lego/colour_hex.go
