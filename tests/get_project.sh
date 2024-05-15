#!/bin/bash

url="http://localhost:8000/"
email="oli.amb5@gmail.com"
pw="mypassword"
endpoint='project/get/'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# Get project
query="Bitcoin"
echo "Endpoint: "$endpoint$query
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint$query
