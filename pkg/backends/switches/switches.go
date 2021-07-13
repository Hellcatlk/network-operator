package switches

import (
	"context"
	"fmt"

	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/backends/switches/ansible"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

type newFuncType func(context.Context, *provider.Config) (backends.Switch, error)

var backendNews map[string]newFuncType = make(map[string]newFuncType)

func init() {
	Register("ansible", ansible.New)
}

// Register switch backend
func Register(backend string, new newFuncType) {
	backendNews[backend] = new
}

// New return a implementation of SwitchBackend interface
func New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	if backendNews[config.Backend] == nil {
		return nil, fmt.Errorf("the type of backend(%s) is invalid", config.Backend)
	}

	return backendNews[config.Backend](ctx, config)
}
