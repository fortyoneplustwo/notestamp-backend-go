#!/bin/bash

url="http://localhost:8000/"
endpoint='/logout'

# grab cookies
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# log out
echo "Endpoint: "$endpoint
curl -X POST -b cookies.txt -v $url$endpoint
echo
