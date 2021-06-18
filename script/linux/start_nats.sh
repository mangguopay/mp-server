#!/bin/sh
nohup nats-streaming-server -m 8222 -p 4222  -DV --cluster_id push_mq > run_nats.log 2>&1 &
