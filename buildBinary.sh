#!/bin/zsh
ROOT_DIR=${PWD}
go build "$ROOT_DIR"/greet/
go build "$ROOT_DIR"/greet/client/
go build "$ROOT_DIR"/greet/server/
