package switches

import (
	"context"
	"testing"

	"github.com/Hellcatlk/network-operator/pkg/provider"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name        string
		OS          string
		Protocol    string
		expectError bool
	}{
		{
			name:        "new test switch",
			OS:          "test",
			Protocol:    "test",
			expectError: false,
		},
		{
			name:        "new not existed switch",
			OS:          "notExisted",
			Protocol:    "test",
			expectError: true,
		},
		{
			name:        "input invalid url",
			OS:          "test",
			Protocol:    "invalid protocol",
			expectError: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(context.Background(), &provider.Config{
				OS:       c.OS,
				Protocol: c.Protocol,
			})
			if (err != nil) != c.expectError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
