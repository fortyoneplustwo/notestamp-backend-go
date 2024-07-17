#!/bin/bash

domain="http://localhost:8000"
endpoint="/auth/login"

# usage
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 email password"
    exit 1
fi

email=$1 
pw=$2

# log in
echo "Endpoint: "$endpoint
curl -X POST -d "username="$email"&password="$pw -c cookies.txt -i $domain$endpoint
echo -e "\n"
