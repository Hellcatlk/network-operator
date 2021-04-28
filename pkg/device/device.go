package device

import (
	"context"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
)

// Switch is a interface for different protocol
type Switch interface {
	// PowerOn
	PowerOn(ctx context.Context) (err error)

	// PowerOff
	PowerOff(ctx context.Context) (err error)

	// GetPortAttr
	GetPortAttr(ctx context.Context, portID string) (configuration *v1alpha1.SwitchPortConfiguration, err error)

	// SetPortAttr
	SetPortAttr(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error)

	// ResetPort
	ResetPort(ctx context.Context, portID string) (err error)
}
