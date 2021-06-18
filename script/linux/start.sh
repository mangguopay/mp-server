#!/bin/bash
cd /home/software/run/go-micro/mp-server/${1} && nohup /home/software/run/go-micro/mp-server/${1}/${1}d > ${1}d.log 2>&1 &
