#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 project_title"
    exit 1
fi

title_encoded=$(jq -rn --arg str "$1" '$str|@uri')

url="http://localhost:8000/"
endpoint='media/download/'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# Get project
echo "Endpoint: "$endpoint$title_encoded
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint$title_encoded \
  --output media.pdf
