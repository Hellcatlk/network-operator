#!/usr/bin/sh

set -ue

os=$(go env GOOS)
arch=$(go env GOARCH)

mkdir -p ./bin
curl -L https://github.com/golangci/golangci-lint/releases/download/v1.38.0/golangci-lint-1.38.0-"${os}"-"${arch}".tar.gz | tar -xz -C ./bin/
mv ./bin/golangci-lint-1.38.0-"${os}"-"${arch}"/golangci-lint ./bin
rm -rf ./bin/golangci-lint-1.38.0-"${os}"-"${arch}"
