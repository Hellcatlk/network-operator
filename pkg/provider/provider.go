package provider

import (
	"context"
	"fmt"

	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Config is configuration of provider device
type Config struct {
	OS       string
	Protocol string
	Host     string
	Cert     *certificate.Certificate
	Options  map[string]string
}

// Switch is a interface of provider switch
type Switch interface {
	// GetOS return switch's os
	GetOS() string

	// GetProtocol return switch's protocol
	GetProtocol() string

	// GetHost return switch's host
	GetHost() string

	// GetSecret return switch's certificate secret reference
	GetSecret() *corev1.SecretReference

	// GetOptions return switch's options
	GetOptions() map[string]string
}

// GetSwitchConfiguration get configuration of provider switch
func GetSwitchConfiguration(ctx context.Context, client client.Client, sw Switch) (*Config, error) {
	if sw == nil {
		return nil, fmt.Errorf("provider switch is nil")
	}

	cert, err := certificate.Fetch(ctx, client, sw.GetSecret())
	if err != nil {
		return nil, err
	}

	return &Config{
		OS:       sw.GetOS(),
		Protocol: sw.GetProtocol(),
		Host:     sw.GetHost(),
		Cert:     cert,
		Options:  sw.GetOptions(),
	}, nil
}
