package switches

import (
	"context"
	"fmt"
	"net/url"

	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device/switches/openvswitch"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device/switches/test"
)

type newType func(ctx context.Context, host string, username string, password string, options map[string]string) (sw device.Switch, err error)

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
func New(ctx context.Context, os string, rawurl string, username string, password string, options map[string]string) (sw device.Switch, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if news[os] == nil {
		return nil, fmt.Errorf("invalid OS %s", os)
	}

	new := news[os][u.Scheme]
	if new == nil {
		return nil, fmt.Errorf("invalid scheme %s", u.Scheme)
	}
	return new(ctx, u.Host, username, password, options)
}
