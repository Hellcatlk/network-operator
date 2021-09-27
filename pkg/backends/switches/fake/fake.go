package fake

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

// New return fake backend
func New(ctx context.Context, config *provider.SwitchConfiguration) (backends.Switch, error) {
	return &fake{}, nil
}

// Test just for test
type fake struct {
}

// IsAvailable check switch is available or not
func (t *fake) IsAvailable() error {
	return nil
}

// GetPortAttr just for test
func (t *fake) GetPortAttr(ctx context.Context, name string) (*v1alpha1.SwitchPortConfigurationSpec, error) {
	return &v1alpha1.SwitchPortConfigurationSpec{}, nil
}

// SetPortAttr just for test
func (t *fake) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfigurationSpec) error {
	return nil
}

// ResetPort just for test
func (t *fake) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfigurationSpec) error {
	return nil
}
