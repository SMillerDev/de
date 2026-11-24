#!/usr/bin/env bash

title=$1

if [[ -z "$title" ]]; then
  echo "No title specified!" >&2
  exit 1
fi

ukr_title=$(curl --fail --silent --get \
    --data-urlencode "titles=$title" \
    --data-urlencode "lllimit=500" \
    --data-urlencode "format=json" \
    --data-urlencode "prop=langlinks" \
    --data-urlencode "action=query" \
    'https://en.wikipedia.org/w/api.php' | jq -r '.query.pages[].langlinks | map(select(.lang == "uk")) | .[0]["*"]')

en_extr=$(curl --fail --silent --get \
    --data-urlencode "titles=$title" \
    --data-urlencode "action=query" \
    --data-urlencode "exsectionformat=plain" \
    --data-urlencode "explaintext=true" \
    --data-urlencode "prop=extracts" \
    --data-urlencode "exintro=true" \
    --data-urlencode "exsentences=2" \
    --data-urlencode "format=json" \
     'https://en.wikipedia.org/w/api.php' | jq '.query.pages[] | {name: .title, prompt: {en: .extract}}')

uk_extr=$(curl --fail --silent --get \
    --data-urlencode "titles=$ukr_title" \
    --data-urlencode "action=query" \
    --data-urlencode "exsectionformat=plain" \
    --data-urlencode "explaintext=true" \
    --data-urlencode "prop=extracts" \
    --data-urlencode "exintro=true" \
    --data-urlencode "exsentences=2" \
    --data-urlencode "format=json" \
     'https://uk.wikipedia.org/w/api.php' | jq '.query.pages[] | {prompt: {uk: .extract}}')

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

query_osm() {
    local title
    title="$1"
    if [[ $title =~ [0-9]+$ ]]; then
        echo "${title} is an id" >&2
        IFS="-" read -r type id <<< "$title"
        print_json "$type" "$id"
        return
    fi

    echo "${title} is a search term" >&2
    data=$(curl --fail --silent  --get --data-urlencode "q=${title}" --data-urlencode "format=json" "https://nominatim.openstreetmap.org/search" | jq)
    if [[ "$data" == "[]" ]]; then
        echo "No data for $title" >&2
        exit 1
    fi
    id=$(echo "$data" | jq -r ".[0].osm_id")
    type=$(echo "$data" | jq -r ".[0].osm_type")
    print_json "$type" "$id"
}

osm=""
osm_data=$(query_osm "$title")
if [[ -n "$osm_data" ]]; then
    osm=$(echo "$osm_data" | jq '{geometry: .}')
fi

echo "$en_extr" "$uk_extr" "$osm" |\
  jq --slurp 'reduce .[] as $item ({}; . * $item)'
