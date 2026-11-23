#!/usr/bin/env bash
query="$1"

print_json() {
    if [[ $1 == "way" ]]; then
        osmexp way "$2"  | jq --compact-output '{type:"Feature",properties:{},geometry: .}'
    elif [[ $1 == "relation" ]]; then
        osmexp rel "$2"  | jq --compact-output .features[0]
    elif [[ $1 == "node" ]]; then
        osmexp node "$2"  | jq --compact-output
    else
        echo "No mapping found for $1 $2" >&2
    fi
}

if [[ $query =~ [0-9]+$ ]]; then
    echo "${query} is an id" >&2
    IFS="-" read -r type id <<< "$query"
    print_json "$type" "$id"
else
    echo "${query} is a search term"
    data=$(curl --fail --silent  --get --data-urlencode "q=${query}" --data-urlencode "format=json" "https://nominatim.openstreetmap.org/search" | jq)
    if [[ "$data" == "[]" ]]; then
       echo "No data" >&2
       exit 1
    fi
    id=$(echo "$data" | jq -r ".[0].osm_id")
    type=$(echo "$data" | jq -r ".[0].osm_type")
    print_json "$type" "$id"
fi
