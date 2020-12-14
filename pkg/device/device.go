package device

import (
	"context"
	"fmt"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PortState ...
type PortState string

const (
	// None ...
	None PortState = "none"

	// Configuring ...
	Configuring PortState = "configuring"

	// Configured ...
	Configured PortState = "configured"

	// ConfigureFailed ...
	ConfigureFailed PortState = "configure failed"

	// Deleting ...
	Deleting PortState = "deleting"

	// Deleted ...
	Deleted PortState = "deleted"

	// DeleteFailed ...
	DeleteFailed PortState = "delete failed"
)

// New ...
func New(ctx context.Context, client *client.Client, deviceRef *v1alpha1.DeviceRef) (device Device, err error) {
	// Deal possible panic
	defer func() {
		err := recover()
		if err != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	switch deviceRef.Kind {
	case "Switch":
		device, err = newSwitch(ctx, client, deviceRef)
	default:
		err = fmt.Errorf("no device for the kind(%s)", deviceRef.Kind)
	}

	return
}

// Device ...
type Device interface {
	// ConfigurePort set the network configure to the port
	ConfigurePort(ctx context.Context, networkConfiguration *v1alpha1.NetworkConfiguration, port *v1alpha1.NetworkBindingSpecPort) error

	// DeConfigurePort remove the network configure from the port
	DeConfigurePort(ctx context.Context, port *v1alpha1.NetworkBindingSpecPort) error

	// PortState return the port's state of the device
	PortState(ctx context.Context, portID string) PortState
}
