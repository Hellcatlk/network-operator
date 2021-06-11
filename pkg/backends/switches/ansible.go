package switches

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

func init() {
	Register("ansible", &ansible{})
}

// Test just for test
type ansible struct {
}

// New return test backend
func (s *ansible) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	return &ansible{}, nil
}

// PowerOn just for test
func (s *ansible) PowerOn(ctx context.Context) error {
	return nil
}

// PowerOff just for test
func (s *ansible) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr just for test
func (s *ansible) GetPortAttr(ctx context.Context, name string) (*v1alpha1.SwitchPortConfiguration, error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr just for test
func (s *ansible) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}

// ResetPort just for test
func (s *ansible) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}
