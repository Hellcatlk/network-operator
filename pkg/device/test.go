package device

import (
	"context"
)

// newSwitch ...
func newTest() (*Test, error) {

	return &Test{}, nil
}

// Test is a kind of network device
type Test struct {
}

// ConfigurePort set the network configure to the port
func (s *Test) ConfigurePort(ctx context.Context, configuration interface{}, portID string) error {
	return nil
}

// DeConfigurePort remove the network configure from the port
func (s *Test) DeConfigurePort(ctx context.Context, portID string) error {
	return nil
}

// CheckPortConfigutation checks whether the configuration is configured on the port
func (s *Test) CheckPortConfigutation(ctx context.Context, configuration interface{}, portID string) (bool, error) {
	return true, nil
}
