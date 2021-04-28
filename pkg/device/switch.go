package device

import (
	"context"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
)

// Switch is a interface
type Switch interface {
	PowerOn(ctx context.Context) (err error)
	PowerOff(ctx context.Context) (err error)
	CreateVlan(ctx context.Context, vlans []v1alpha1.VLAN) (err error)
	DeleteVlan(ctx context.Context, vlans []v1alpha1.VLAN) (err error)
	GetPortAttr(ctx context.Context, portID string) (vlans []v1alpha1.VLAN, portType v1alpha1.PortType, err error)
	SetPortAttr(ctx context.Context, portID string, vlans []v1alpha1.VLAN, portType v1alpha1.PortType) (err error)
	ResetPort(ctx context.Context, portID string) (err error)
}
