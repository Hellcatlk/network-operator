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

// NewCLI return openvswitch cli backend, username and password is invalid for the backend

func NewCLI(ctx context.Context, address string, username string, password string, options map[string]string) (sw device.Switch, err error) {
	if len(options) == 0 || options["bridge"] == "" {
		return nil, fmt.Errorf("bridge of openvswitch cli backend is required")
	}

	return &CLI{
		address:  address,
		username: username,
		password: password,
		bridge:   options["bridge"],
	}, nil
}

// CLI control local openvswitch by CLI
type CLI struct {
	address  string
	username string
	password string
	bridge   string
}

// PowerOn enable openvswitch
func (c *CLI) PowerOn(ctx context.Context) (err error) {
	return ssh.Run(c.address, c.username, c.password, exec.Command("ovs-vsctl", "list", "br", c.bridge)) // #nosec
}

// PowerOff disable openvswitch
func (c *CLI) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr get the port's configure
func (c *CLI) GetPortAttr(ctx context.Context, portID string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr set configure to the port
func (c *CLI) SetPortAttr(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	var tag string = ""
	for _, vlan := range configuration.Spec.Vlans {
		if tag == "" {
			tag = "tag=" + strconv.Itoa(int(vlan.ID))
		} else {
			tag = tag + "," + strconv.Itoa(int(vlan.ID))
		}
	}

	return ssh.Run(c.address, c.username, c.password, exec.Command("ovs-vsctl", "set ", "port", portID, tag)) // #nosec
}

// ResetPort remove all configure of the port
func (c *CLI) ResetPort(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	return ssh.Run(c.address, c.username, c.password, exec.Command("ovs-vsctl", "set ", "port", portID, "tag=")) // #nosec
}
