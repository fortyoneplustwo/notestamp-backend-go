#!/bin/bash

url="http://localhost:8000/"
email="oli.amb5@gmail.com"
pw="mypassword"
endpoint='project/delete/'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# delete project
query="my%20project%203"
echo "Endpoint: "$endpoint$query
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint$query
echo "-----------------------------------------------------"

