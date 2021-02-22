package device

import (
	"context"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// newSwitch ...
func newSwitch(ctx context.Context, client client.Client, deviceRef *metav1.OwnerReference) (*Switch, error) {
	instance := &v1alpha1.Switch{}

	// Get SwitchDevice CR
	err := client.Get(
		ctx,
		types.NamespacedName{
			Name: deviceRef.Name,
		},
		instance,
	)
	if err != nil {
		return nil, err
	}

	return &Switch{
		client:   client,
		instance: instance,
	}, nil
}

// Switch is a kind of network device
type Switch struct {
	client   client.Client
	instance *v1alpha1.Switch
}

// ConfigurePort set the network configure to the port
func (s *Switch) ConfigurePort(ctx context.Context, configuration interface{}, portID string) error {
	return nil
}

// DeConfigurePort remove the network configure from the port
func (s *Switch) DeConfigurePort(ctx context.Context, portID string) error {
	return nil
}

// CheckPortConfigutation checks whether the configuration is configured on the port
func (s *Switch) CheckPortConfigutation(ctx context.Context, configuration interface{}, portID string) bool {
	return false
}
