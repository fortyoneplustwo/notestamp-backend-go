#!/bin/bash

url="http://localhost:8000/"
email="oli.amb5@gmail.com"
pw="mypassword"
endpoint='project/list'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# list projects
echo "Endpoint: "$endpoint
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint
echo -e "\n"
