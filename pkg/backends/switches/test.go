package switches

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

func init() {
	Register("test", &test{})
}

// Test just for test
type test struct {
}

// New return test backend
func (t *test) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	return &test{}, nil
}

// PowerOn just for test
func (t *test) PowerOn(ctx context.Context) error {
	return nil
}

// PowerOff just for test
func (t *test) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr just for test
func (t *test) GetPortAttr(ctx context.Context, name string) (*v1alpha1.SwitchPortConfiguration, error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr just for test
func (t *test) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}

// ResetPort just for test
func (t *test) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}
