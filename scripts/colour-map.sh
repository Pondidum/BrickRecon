#!/bin/sh

url="https://api.brickowl.com/v1/catalog/color_list?key=${BRICKOWL_API_KEY}"
buildMap='. | map({ (.bl_ids[0]): .hex }) | add'
convert_keys='s/"\(.*\)":/\1:/g'
add_trailing_comma='s/\(.*\)"/\1",/'

json=$(curl "$url" --silent | jq --tab --sort-keys "$buildMap" | sed "$convert_keys" | sed -z "$add_trailing_comma")

printf "package lego\n\nvar hexColours = map[int]string%s" "$json" > lego/colour_hex.go
