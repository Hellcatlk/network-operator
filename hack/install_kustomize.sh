#!/bin/sh

set -ue

os=$(go env GOOS)
arch=$(go env GOARCH)

mkdir -p ./bin
curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv4.0.5/kustomize_v4.0.5_"${os}"_"${arch}".tar.gz | tar -xz -C ./bin/
