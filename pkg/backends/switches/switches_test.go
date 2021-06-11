package switches

import (
	"context"
	"testing"

	"github.com/Hellcatlk/network-operator/pkg/provider"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name        string
		backend     string
		os          string
		expectError bool
	}{
		{
			name:        "new not existed switch",
			backend:     "notExisted",
			expectError: true,
		},
		{
			name:        "new test switch",
			backend:     "test",
			expectError: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(context.Background(), c.backend, &provider.Config{})
			if (err != nil) != c.expectError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
