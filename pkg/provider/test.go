package provider

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TestSwitch is a instance of provider switch
type TestSwitch struct {
}

// GetConfiguration generate configuration from provider switch
func (t *TestSwitch) GetConfiguration(ctx context.Context, client client.Client) (*SwitchConfiguration, error) {
	return &SwitchConfiguration{
		Backend: "test",
	}, nil
}
