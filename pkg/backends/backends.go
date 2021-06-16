package backends

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/provider"
)

// Switch is a interface for switch backend
type Switch interface {
	// New return backend itself
	New(ctx context.Context, config *provider.Config) (Switch, error)

	// PowerOn enable switch
	PowerOn(ctx context.Context) error

	// PowerOff disable switch
	PowerOff(ctx context.Context) error

	// GetPortAttr get the port's configure
	GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfiguration, error)

	// SetPortAttr set configure to the port
	SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error

	// ResetPort remove all configure of the port
	ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error
}
