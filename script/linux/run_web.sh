#!/bin/sh
nohup micro --registry=etcd web --namespace=go.micro.web > go_web.log 2>&1 &
