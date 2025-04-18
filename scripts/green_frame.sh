#!/bin/bash

curl -X POST http://localhost:17000/ -d $'green\nbgrect 0.25 0.25 0.75 0.75\nupdate'
