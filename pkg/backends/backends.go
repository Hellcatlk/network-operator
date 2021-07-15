package backends

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
)

// Switch is a interface for switch backend
type Switch interface {
	// GetPortAttr get the port's configuration
	GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfiguration, error)

	// SetPortAttr set configure to the port
	SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error

	// ResetPort remove all configure of the port
	ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error
}
