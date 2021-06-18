#!/bin/bash
echo 'build '${1}
cd /home/software/run/go-micro/mp-server/$1 && go build -o ${1}d . 
