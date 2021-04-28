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

	// CreateVlan
	CreateVlan(ctx context.Context, vlans []v1alpha1.VLAN) (err error)

	// DeleteVlan
	DeleteVlan(ctx context.Context, vlans []v1alpha1.VLAN) (err error)

	// GetPortAttr
	GetPortAttr(ctx context.Context, portID string) (vlans []v1alpha1.VLAN, portType v1alpha1.PortType, err error)

	// SetPortAttr
	SetPortAttr(ctx context.Context, portID string, vlans []v1alpha1.VLAN, portType v1alpha1.PortType) (err error)

	// ResetPort
	ResetPort(ctx context.Context, portID string) (err error)
}
