#!/bin/bash
# Generate Go code from the .proto file
protoc --go_out=./gen --go-grpc_out=./gen proto/library/v1/*.proto