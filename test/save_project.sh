#!/bin/bash

url="http://localhost:8000/"
endpoint="project/save"

# usage
if [ "$#" -ne 7 ]; then
  echo "Usage: $0 notes_file media_file type title label src mimetype"
    exit 1
fi

notes=$1
media=$2
type=$3
title=$4
label=$5
src=$6
mimetype=$7

json=$(jq -n --arg type "$type" --arg title "$title"  --arg label "$label"  --arg src "$src"  --arg mimetype "$mimetype" '$ARGS.named')

metadata='metadata='$json

# grab cookies from cookies.txt
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

if [ "$media" != "" ]; then
  # save project with media
  echo "Endpoint: "$endpoint
  curl -X POST\
    -H "Content-Type: multipart/form-data" \
    -b "refresh-token="$rtk";access-token="$atk \
    -c cookies.txt \
    -F "notesFile=@"$notes  \
    -F "mediaFile=@"$media  \
    -F "$metadata" \
    -i \
    $url$endpoint
  echo -e "\n"
else
  # Save without media
  echo "Endpoint: "$endpoint
  curl -X POST\
    -H "Content-Type: multipart/form-data" \
    -b "refresh-token="$rtk";access-token="$atk \
    -c cookies.txt \
    -F "notesFile=@"$notes  \
    -F "$metadata" \
    -i \
    $url$endpoint
  echo -e "\n"
fi


