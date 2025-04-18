#!/bin/bash

curl -X POST http://localhost:17000/ -d $'reset\nwhite\nfigure 0.50 0.50\nupdate'
sleep 1
curl -X POST http://localhost:17000/ -d $'reset\nwhite\nfigure 0.53 0.53\nupdate'
sleep 1
curl -X POST http://localhost:17000/ -d $'reset\nwhite\nfigure 0.56 0.56\nupdate'
sleep 1
curl -X POST http://localhost:17000/ -d $'reset\nwhite\nfigure 0.59 0.59\nupdate'
sleep 1
curl -X POST http://localhost:17000/ -d $'reset\nwhite\nfigure 0.62 0.62\nupdate'
