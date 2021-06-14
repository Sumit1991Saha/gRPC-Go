#!/bin/zsh

## old style to generate grpc code
##protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.

# to generate grpc code
protoc blog/blogpb/blog.proto --go-grpc_out=.

## to generate proto messages
protoc blog/blogpb/blog.proto --go_out=.

