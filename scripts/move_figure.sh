#!/bin/bash

positions=("0.50" "0.53" "0.56" "0.59" "0.62")

for pos in "${positions[@]}"
do
    curl -X POST http://localhost:17000/ -d $'reset\nwhite\nfigure '"$pos $pos"$'\nupdate'
    sleep 1
done