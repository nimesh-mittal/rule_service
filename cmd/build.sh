#!/usr/bin/env bash

# generate swagger document
swag init -d /opt/dev/code/ws-go/src/location_service/resources/ -g /opt/dev/code/ws-go/src/location_service/main.go

# create docker image
docker build --tag location_service_v4 .