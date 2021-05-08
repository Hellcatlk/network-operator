package openvswitch

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
	"github.com/metal3-io/networkconfiguration-operator/pkg/utils/ssh"
)

// NewSSH return openvswitch ssh backend
func NewSSH(ctx context.Context, address string, username string, password string, options map[string]string) (sw device.Switch, err error) {
	if len(options) == 0 || options["bridge"] == "" {
		return nil, fmt.Errorf("bridge of openvswitch cli backend is required")
	}

	return &SSH{
		address:  address,
		username: username,
		password: password,
		bridge:   options["bridge"],
	}, nil
}

// CLI control local openvswitch by CLI
type SSH struct {
	address  string
	username string
	password string
	bridge   string
}

// PowerOn enable openvswitch
func (c *SSH) PowerOn(ctx context.Context) (err error) {
	return ssh.Run(c.address, c.username, c.password, exec.Command("ovs-vsctl", "list", "br", c.bridge)) // #nosec
}

// PowerOff disable openvswitch
func (c *SSH) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr get the port's configure
func (c *SSH) GetPortAttr(ctx context.Context, portID string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr set configure to the port
func (c *SSH) SetPortAttr(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	if len(configuration.Spec.Vlans) != 1 {
		return fmt.Errorf("vlans's len port of openvswitch is 1")
	}

	return ssh.Run(c.address, c.username, c.password, exec.Command("ovs-vsctl", "set", "port", portID, "tag="+strconv.Itoa(int(configuration.Spec.Vlans[0].ID)))) // #nosec
}

// ResetPort remove all configure of the port
func (c *SSH) ResetPort(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	if len(configuration.Spec.Vlans) != 1 {
		return fmt.Errorf("vlans's len port of openvswitch is 1")
	}

	return ssh.Run(c.address, c.username, c.password, exec.Command("ovs-vsctl", "remove ", "port", portID, "tag", strconv.Itoa(int(configuration.Spec.Vlans[0].ID)))) // #nosec
}
