#!/bin/zsh

## old style to generate grpc code
##protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc greet/greetpb/greet.proto --go-grpc_out=.