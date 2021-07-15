// Package provider convert the content of ProviderSwitch to configuration
package provider

import (
	"context"

	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SwitchConfiguration is configuration of device provider
type SwitchConfiguration struct {
	// which os this provider switch used
	OS string
	// Switch's host
	Host string
	// Certificate of switch
	Cert *certificate.Certificate
	// Which backend to use
	Backend string
	Options map[string]interface{}
}

// Switch is a interface of provider switch
type Switch interface {
	// GetConfiguration generate configuration from provider switch
	GetConfiguration(ctx context.Context, client client.Client) (*SwitchConfiguration, error)
}
