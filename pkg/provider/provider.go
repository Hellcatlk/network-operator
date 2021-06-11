package provider

import (
	"context"

	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Config is configuration of provider device
type Config struct {
	OS      string
	Host    string
	Cert    *certificate.Certificate
	Options map[string]string
}

// Switch is a interface of provider switch
type Switch interface {
	// GetConfiguration generate configuration from provider switch
	GetConfiguration(ctx context.Context, client client.Client) (*Config, error)
}
