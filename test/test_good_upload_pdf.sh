#!/bin/bash

email="oli.amb5@gmail.com"
pw="mypassword"
sep="==================================================="

./login.sh "$email" "$pw"
echo "$sep"
./list_projects.sh
echo "$sep"

notes="document.stmp"
media="bitcoin_satoshi.pdf"
type="pdf"
title="Pdf project"
label="PDF Reader"
src=""
mimetype="application/pdf"

# Upload with notes and media
./save_project.sh "$notes" "$media" "$type" "$title" "$label" "$src" "$mimetype" 
echo "$sep"
./delete_project.sh "$title"
echo "$sep"

# Uplaod without media
title="Pdf project (empty)"
media=""
./save_project.sh "$notes" "$media" "$type" "$title" "$label" "$src" "$mimetype" 
echo "$sep"
./delete_project.sh "$title"
echo "$sep"

./logout.sh

