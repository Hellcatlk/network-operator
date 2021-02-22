package device

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New ...
func New(ctx context.Context, client client.Client, deviceRef *metav1.OwnerReference) (device Device, err error) {
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
	case "Test":
		device, err = newTest()
	default:
		err = fmt.Errorf("no device for the kind(%s)", deviceRef.Kind)
	}

	return
}

// Device ...
type Device interface {
	// ConfigurePort set the network configure to the port
	ConfigurePort(ctx context.Context, configuration interface{}, portID string) error

	// DeConfigurePort remove the network configure from the port
	DeConfigurePort(ctx context.Context, portID string) error

	// CheckPortConfigutation checks whether the configuration is configured on the port
	CheckPortConfigutation(ctx context.Context, configuration interface{}, portID string) (bool, error)
}
