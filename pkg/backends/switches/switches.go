package switches

import (
	"context"
	"fmt"

	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

var switchBackends map[string]backends.Switch = make(map[string]backends.Switch)

// Register switch backend
func Register(backendType string, backend backends.Switch) {
	switchBackends[backendType] = backend
}

// New return a implementation of SwitchBackend interface
func New(ctx context.Context, backendType string, config *provider.Config) (backends.Switch, error) {
	if switchBackends[backendType] == nil {
		return nil, fmt.Errorf("the type of backend(%s) is invalid", backendType)
	}

	return switchBackends[backendType].New(ctx, config)
}
