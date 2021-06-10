package devices

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
)

// Switch is a interface for different protocol
type Switch interface {
	// PowerOn enable switch
	PowerOn(ctx context.Context) error

	// PowerOff disable switch
	PowerOff(ctx context.Context) error

	// GetPortAttr get the port's configure
	GetPortAttr(ctx context.Context, name string) (*v1alpha1.SwitchPortConfiguration, error)

	// SetPortAttr set configure to the port
	SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error

	// ResetPort remove all configure of the port
	ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error
}
