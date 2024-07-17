#!/bin/bash

# usage
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 email password"
    exit 1
fi

email=$1 
pw=$2

url="http://localhost:8000/auth"

# register
endpoint="/register"
echo "Endpoint: /register"
curl -X POST -d "username="$email"&password="$pw -s -i $url$endpoint
echo -e "\n"

