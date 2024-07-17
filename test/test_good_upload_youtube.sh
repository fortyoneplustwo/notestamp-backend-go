#!/bin/bash

email="oli.amb5@gmail.com"
pw="mypassword"
sep="==================================================="

./login.sh "$email" "$pw"
echo "$sep"
./list_projects.sh
echo "$sep"

notes="document.stmp"
media=""
type="pdf"
title="Youtube project"
label="Youtube Player"
src="https://www.youtube.com/watch?v=SLosyUI5puk"
mimetype=""

# Upload with notes and media
./save_project.sh "$notes" "$media" "$type" "$title" "$label" "$src" "$mimetype" 
echo "$sep"
./delete_project.sh "$title"
echo "$sep"

# Upload without src
src=""
./save_project.sh "$notes" "$media" "$type" "$title" "$label" "$src" "$mimetype" 
echo "$sep"
./delete_project.sh "$title"
echo "$sep"

./logout.sh
