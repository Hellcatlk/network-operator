package test

import (
	"context"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
)

// NewTest return test backend
func NewTest(ctx context.Context, host string, username string, password string) (sw device.Switch, err error) {
	return &Test{}, nil
}

// Test just for test
type Test struct {
}

// PowerOn just for test
func (t *Test) PowerOn(ctx context.Context) (err error) {
	return
}

// PowerOff just for test
func (t *Test) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr just for test
func (t *Test) GetPortAttr(ctx context.Context, portID string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr just for test
func (t *Test) SetPortAttr(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	return
}

// ResetPort just for test
func (t *Test) ResetPort(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	return
}
