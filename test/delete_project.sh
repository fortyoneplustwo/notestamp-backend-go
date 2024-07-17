#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 project_title"
    exit 1
fi

title_encoded=$(jq -rn --arg str "$1" '$str|@uri')

url="http://localhost:8000/"
endpoint='project/delete/'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# delete project
echo "Endpoint: "$endpoint$title_encoded
curl -X DELETE \
  -b "refresh-token="$rtk";access-token="$atk \
  -c cookies.txt \
  -i $url$endpoint$title_encoded
echo -e "\n"

