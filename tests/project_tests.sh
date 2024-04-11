#!/bin/bash

url="http://localhost:8000/"
email="oli.amb5@gmail.com"
pw="mypassword"

# log in
endpoint="auth/login"
echo "Endpoint: "$endpoint
curl -X POST -d "username=oli.amb5@gmail.com&password=mypassword" -c cookies.txt -i $url$endpoint
echo -e "\n"

atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# save project
endpoint="project/save"
echo "Endpoint: "$endpoint
curl -X POST \
  -b "refresh-token="$rtk";access-token="$atk \
  -c cookies.txt \
  -F "content=notes"  \
  -F 'metadata={"type": "audio", "title": "my project 3", "label": "Audio Player", "src": "", "mimetype": ""}' \
  -i \
  $url$endpoint
echo -e "\n"

atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# list projects
endpoint='project/list'
echo
echo "Endpoint: "$endpoint
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint
echo

atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# Get project
endpoint='project/get/'
query="my%20project%203"
echo
echo "Endpoint: "$endpoint$query
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint$query
echo

atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# Delete project
endpoint='project/delete/'
query="my%20project%203"
echo
echo "Endpoint: "$endpoint$query
curl -X GET \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint$query
echo

# log out
endpoint='auth/logout'
echo
echo "Endpoint: "$endpoint
curl -X POST \
  -b "refresh-token="$rtk";access-token="$atk \
  -i $url$endpoint
echo

