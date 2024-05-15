#!/bin/bash

url="http://localhost:8000/"
email="oli.amb5@gmail.com"
pw="mypassword"
endpoint="project/save"

# grab cookies from cookies.txt
atk=$(grep 'access-token' cookies.txt | awk '{print $NF}')
rtk=$(grep 'refresh-token' cookies.txt | awk '{print $NF}')

# save project
echo "Endpoint: "$endpoint
curl -X POST \
  -b "refresh-token="$rtk";access-token="$atk \
  -c cookies.txt \
  -F "notesFile=@document.stmp"  \
  -F "mediaFile=@bitcoin_satoshi.pdf"  \
  -F 'metadata={"type": "pdf", "title": "Bitcoin", "label": "PDF Reader", "src": "", "mimetype": "application/pdf"}' \
  -i \
  $url$endpoint
echo -e "\n"

