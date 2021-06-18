#!/bin/bash
cd common && make && cd ..
sh build.sh api-mobile
sh build.sh api-webadmin
sh build.sh sms-srv
sh build.sh auth-srv
sh build.sh cust-srv
sh build.sh bill-srv
