#!/bin/sh
nohup micro  --registry=etcd api --handler=proxy --namespace=go.micro.api --address=0.0.0.0:8888 > api.log 2>&1 &
