package openvswitch

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
	"github.com/metal3-io/networkconfiguration-operator/pkg/utils/certificate"
	"github.com/metal3-io/networkconfiguration-operator/pkg/utils/ssh"
)

// NewSSH return openvswitch ssh backend
func NewSSH(ctx context.Context, Host string, cert *certificate.Certificate, options map[string]string) (sw device.Switch, err error) {
	if len(options) == 0 || options["bridge"] == "" {
		return nil, fmt.Errorf("bridge of openvswitch cli backend is required")
	}

	return &SSH{
		Host:     Host,
		username: cert.Username,
		password: cert.Password,
		bridge:   options["bridge"],
	}, nil
}

// SSH control openvswitch by ssh and cli
type SSH struct {
	Host     string
	username string
	password string
	bridge   string
}

// PowerOn enable openvswitch
func (c *SSH) PowerOn(ctx context.Context) (err error) {
	err = ssh.Run(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "list", "br", c.bridge,
	)) // #nosec
	if err != nil {
		return fmt.Errorf("check birdge failed: %v", err)
	}

	return nil
}

// PowerOff disable openvswitch
func (c *SSH) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr get the port's configure
func (c *SSH) GetPortAttr(ctx context.Context, name string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	output, err := ssh.Output(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "list", "port", name,
		"|", "grep", "-E", "-w", "^tag",
		"|", "grep", "-o", "[0-9]*",
	)) // #nosec
	if err != nil {
		return nil, fmt.Errorf("get port failed: %v", err)
	}

	id, err := strconv.Atoi(strings.Trim(string(output), "\n"))
	if err != nil {
		return &v1alpha1.SwitchPortConfiguration{}, nil
	}

	return &v1alpha1.SwitchPortConfiguration{
		Spec: v1alpha1.SwitchPortConfigurationSpec{
			Vlans: []v1alpha1.VLAN{
				{
					ID: id,
				},
			},
		},
	}, nil
}

// SetPortAttr set configure to the port
func (c *SSH) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	if len(configuration.Spec.Vlans) != 1 {
		return fmt.Errorf("vlans's len port of openvswitch is 1")
	}

	err = ssh.Run(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "set", "port", name, "tag="+strconv.Itoa(int(configuration.Spec.Vlans[0].ID)),
	)) // #nosec
	if err != nil {
		return fmt.Errorf("set port failed: %v", err)
	}

	return nil
}

// ResetPort remove all configure of the port
func (c *SSH) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	if len(configuration.Spec.Vlans) != 1 {
		return fmt.Errorf("vlans's len port of openvswitch is 1")
	}

	err = ssh.Run(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "remove ", "port", name, "tag", strconv.Itoa(int(configuration.Spec.Vlans[0].ID)),
	)) // #nosec
	if err != nil {
		return fmt.Errorf("set port failed: %v", err)
	}

	return nil
}
