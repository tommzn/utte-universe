#!/bin/bash

protoc --go_out=. --go-grpc_out=./ core/proto/game.proto