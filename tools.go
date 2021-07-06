// +build tools

package tools

import (
	_ "github.com/securego/gosec/cmd/gosec"
	_ "golang.org/x/lint/golint"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	_ "sigs.k8s.io/kustomize/kustomize/v3"
)
