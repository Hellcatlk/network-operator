package provider

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FakeSwitch is a instance of provider switch
type FakeSwitch struct {
}

// GetConfiguration generate configuration from provider switch
func (t *FakeSwitch) GetConfiguration(ctx context.Context, client client.Client) (*SwitchConfiguration, error) {
	return &SwitchConfiguration{
		Backend: "fake",
	}, nil
}
