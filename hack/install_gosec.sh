#!/usr/bin/sh

set -ue

os=$(go env GOOS)
arch=$(go env GOARCH)

mkdir -p ./bin
curl -L https://github.com/securego/gosec/releases/download/v2.7.0/gosec_2.7.0_"${os}"_"$arch".tar.gz | tar -xz -C ./bin/
rm ./bin/LICENSE.txt
rm ./bin/README.md
