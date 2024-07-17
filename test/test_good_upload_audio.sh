#!/bin/bash

email="oli.amb5@gmail.com"
pw="mypassword"
sep="==================================================="

./login.sh "$email" "$pw"
echo "$sep"
./list_projects.sh
echo "$sep"

notes="document.stmp"
media="detonation.mp3"
type="audio"
title="Audio Project"
label="Audio Player"
src=""
mimetype="audio/mp3"

# Upload with notes and media
./save_project.sh "$notes" "$media" "$type" "$title" "$label" "$src" "$mimetype" 
echo "$sep"
./delete_project.sh "$title"
echo "$sep"

# Uplaod without media
title="Audio Project (empty)"
media=""
./save_project.sh "$notes" "$media" "$type" "$title" "$label" "$src" "$mimetype" 
echo "$sep"
./delete_project.sh "$title"
echo "$sep"

./logout.sh

