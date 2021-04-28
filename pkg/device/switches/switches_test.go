package switches

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name        string
		OS          string
		URL         string
		expectError bool
	}{
		{
			name:        "new test switch",
			OS:          "test",
			URL:         "test://1234",
			expectError: false,
		},
		{
			name:        "new not existed switch",
			OS:          "notExisted",
			URL:         "notExisted://1234",
			expectError: true,
		},
		{
			name:        "input invalid url",
			OS:          "test",
			URL:         "invalid url",
			expectError: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(context.Background(), c.OS, c.URL, "", "")
			if (err != nil) != c.expectError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
