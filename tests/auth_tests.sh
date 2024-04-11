#!/bin/bash

url="http://localhost:8000/auth"

# register
endpoint="/register"
echo "Endpoint: /register"
curl -X POST -d "username=oli.amb5@gmail.com&password=mypassword" -s -i $url$endpoint
echo -e "\n"

# log in
endpoint="/login"
echo "Endpoint: "$endpoint
curl -X POST -d "username=oli.amb5@gmail.com&password=mypassword" -c cookies.txt -i $url$endpoint
echo -e "\n"

# log out
endpoint='/logout'
echo "Endpoint: "$endpoint
curl -X POST -b cookies.txt -v $url$endpoint
echo

