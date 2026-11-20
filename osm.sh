#!/usr/bin/env bash
query="$1"
country="${2:-uk}"

data=$(curl --fail --silent  --get --data-urlencode "q=${query}" --data-urlencode "countrycodes=$country" --data-urlencode "format=json" "https://nominatim.openstreetmap.org/search" | jq)
# echo "$data"
if [[ "$data" == "[]" ]]; then
    echo "No data"
    exit 1
fi
id=$(echo "$data" | jq -r ".[0].osm_id")
type=$(echo "$data" | jq -r ".[0].osm_type")

if [[ $type == "way" ]]; then
    osmexp way "$id"  | jq --compact-output '{type:"Feature",properties:{},geometry: .}'
elif [[ $type == "relation" ]]; then
    osmexp rel "$id"  | jq --compact-output .features[0]
else
    echo "No mapping found for $data"
fi
