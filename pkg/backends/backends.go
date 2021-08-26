package backends

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
)

// Switch is a interface for switch backend
type Switch interface {
	// IsAvaliable check switch is avaliable or not
	IsAvaliable() error

	// GetPortAttr get the port's configuration
	GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfigurationSpec, error)

	// SetPortAttr set configure to the port
	SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfigurationSpec) error

	// ResetPort remove all configure of the port
	ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfigurationSpec) error
}
