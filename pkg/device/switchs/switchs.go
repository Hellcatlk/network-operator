package switchs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device/switchs/test"
)

type newType func(ctx context.Context, address string) (sw device.Switch, err error)

var news map[string]map[string]newType

func init() {
	news = make(map[string]map[string]newType, 0)

	// Register switch
	Register("test", "test", test.NewTT)
}

// Register New() function of a switch interface's implementation
func Register(os string, protocolType string, new newType) {
	if news[os] == nil {
		news[os] = make(map[string]newType, 0)
	}

	news[os][protocolType] = new
}

// New return a implementation of switch interface
func New(ctx context.Context, os string, rawurl string) (sw device.Switch, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if news[os] == nil {
		return nil, fmt.Errorf("haven't %s switch type", os)
	}

	new := news[os][u.Scheme]
	if new == nil {
		return nil, fmt.Errorf("invalid scheme %s", u.Scheme)
	}
	return new(ctx, u.Host)
}
