#!/bin/bash

url="http://localhost:8000/auth"

# register
endpoint="/register"
echo "Endpoint: /register"
curl -X POST -d "username=oli.amb5@gmail.com&password=mypassword" -s -i $url$endpoint
echo -e "\n"

