package switches

import (
	"context"
	"fmt"

	"github.com/Hellcatlk/networkconfiguration-operator/pkg/device"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/device/switches/openvswitch"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/device/switches/test"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/provider"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/utils/certificate"
)

type newType func(ctx context.Context, Host string, cert *certificate.Certificate, options map[string]string) (sw device.Switch, err error)

var news map[string]map[string]newType

func init() {
	news = make(map[string]map[string]newType)

	// Register backend
	Register("test", "test", test.NewTest)
	Register("openvswitch", "ssh", openvswitch.NewSSH)
}

// Register New() function of a switch interface's implementation
func Register(os string, protocolType string, new newType) {
	if news[os] == nil {
		news[os] = make(map[string]newType)
	}

	news[os][protocolType] = new
}

// New return a implementation of switch interface
func New(ctx context.Context, config *provider.Config) (sw device.Switch, err error) {

	if news[config.OS] == nil {
		return nil, fmt.Errorf("invalid OS %s", config.OS)
	}

	new := news[config.OS][config.Protocol]
	if new == nil {
		return nil, fmt.Errorf("invalid protocol %s", config.Protocol)
	}

	return new(ctx, config.Host, config.Cert, config.Options)
}
