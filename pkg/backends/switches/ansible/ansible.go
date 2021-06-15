package ansible

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

// Ansible backend
type Ansible struct {
}

// New return test backend
func (a *Ansible) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	return &Ansible{}, nil
}

// PowerOn just for test
func (a *Ansible) PowerOn(ctx context.Context) error {
	return nil
}

// PowerOff just for test
func (a *Ansible) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr just for test
func (a *Ansible) GetPortAttr(ctx context.Context, name string) (*v1alpha1.SwitchPortConfiguration, error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr just for test
func (a *Ansible) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}

// ResetPort just for test
func (a *Ansible) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}
