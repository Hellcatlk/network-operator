package openvswitch

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
)

// NewCLI return openvswitch cli backend, username and password is invalid for the backend
func NewCLI(ctx context.Context, address string, username string, password string) (sw device.Switch, err error) {
	if address == "" {
		return nil, fmt.Errorf("the format of url must be: \"cli:<bridge name>\"")
	}
	return &CLI{
		bridge: address,
	}, nil
}

// CLI control local openvswitch by CLI
type CLI struct {
	bridge string
}

// PowerOn enable openvswitch
func (c *CLI) PowerOn(ctx context.Context) (err error) {
	cmd := exec.Command("ovs-vsctl", "list", "br", c.bridge)
	err = cmd.Run()
	return err
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

	cmd := exec.Command("ovs-vsctl", "set ", "port", portID, tag)
	err = cmd.Run()
	return err
}

// ResetPort remove all configure of the port
func (c *CLI) ResetPort(ctx context.Context, portID string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	cmd := exec.Command("ovs-vsctl", "set ", "port", portID, "tag=")
	err = cmd.Run()
	return err
}
