#!/bin/bash

curl -X POST http://localhost:17000/ -d $'white\nbgrect 0.25 0.25 0.75 0.75\ngreen\nupdate'
