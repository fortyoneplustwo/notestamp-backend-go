#!/bin/bash

# usage
if [ "$#" -ne 0 ]; then
    echo "Usage: $0"
    exit 1
fi

url="http://localhost:8000/"
endpoint='auth/logout'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# log out
echo "Endpoint: "$endpoint
curl -X POST  \
  -b "refresh-token="$rtk";access-token="$atk \
  -c cookies.txt -i $url$endpoint
echo -e "\n"
