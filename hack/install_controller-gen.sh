#!/bin/sh

set -ue

mkdir -p ./bin
curl -L https://github.com/kubernetes-sigs/controller-tools/archive/v0.5.0.tar.gz | tar -xz -C ./bin/
cd bin/controller-tools-0.5.0
go build -o ../controller-gen cmd/controller-gen/main.go
rm -rf ../controller-tools-0.5.0
