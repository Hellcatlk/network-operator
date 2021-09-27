package switches

import (
	"context"
	"testing"

	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/credentials"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name        string
		backend     string
		os          string
		expectError bool
	}{
		{
			name:        "new not existed backend",
			backend:     "notExisted",
			expectError: true,
		},
		{
			name:        "new fake backend",
			backend:     "fake",
			expectError: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(context.Background(), &provider.SwitchConfiguration{
				OS:      "openvswitch",
				Host:    "test",
				Backend: c.backend,
				Credentials: &credentials.Credentials{
					Username: "test",
					Password: "test",
				},
				Options: map[string]interface{}{
					"bridge": "test",
				},
			})
			if (err != nil) != c.expectError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
