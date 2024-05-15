#!/bin/bash

domain="http://localhost:8000"
endpoint="/auth/login"

# log in
echo "Endpoint: "$endpoint
curl -X POST -d "username=oli.amb5@gmail.com&password=mypassword" -c cookies.txt -i $domain$endpoint
echo -e "\n"
