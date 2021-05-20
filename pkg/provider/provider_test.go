package provider

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type fakeClient struct {
	client.Reader
	client.Writer
	client.StatusClient
}

func (c *fakeClient) Get(ctx context.Context, key types.NamespacedName, obj runtime.Object) error {
	return nil
}

func TestFromSwitch(t *testing.T) {
	cases := []struct {
		name           string
		providerSwitch Switch
		expectedError  bool
	}{
		{
			name:           "normal",
			providerSwitch: &TestSwitch{},
		},
		{
			name:           "provider switch is nil",
			providerSwitch: nil,
			expectedError:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := GetSwitchConfiguration(context.Background(), &fakeClient{}, c.providerSwitch)
			if (err != nil) != c.expectedError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
