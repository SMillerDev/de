#!/usr/bin/env bash
query="$1"

id=$(curl --fail --silent  --get --data-urlencode "q=${query}" --data-urlencode "format=json" "https://nominatim.openstreetmap.org/search" | jq -r ".[0].osm_id")
osmexp rel "$id"  | jq --compact-output .features[0]
