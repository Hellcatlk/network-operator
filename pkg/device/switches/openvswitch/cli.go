package openvswitch

import (
	"context"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
)

// NewCLI return openvswitch cli backend
func NewCLI(ctx context.Context, address string, username string, password string) (sw device.Switch, err error) {
	return &CLI{}, nil
}

// CLI control local openvswitch by CLI
type CLI struct {
}

// PowerOn enable openvswitch
func (c *CLI) PowerOn(ctx context.Context) (err error) {
	return
}

// PowerOff disable openvswitch
func (c *CLI) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr get the port's configure
func (c *CLI) GetPortAttr(ctx context.Context, portID string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr set configure to the port
func (c *CLI) SetPortAttr(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	return
}

// ResetPort remove all configure of the port
func (c *CLI) ResetPort(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	return
}
