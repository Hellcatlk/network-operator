package openvswitch

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/Hellcatlk/networkconfiguration-operator/api/v1alpha1"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/devices"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/utils/certificate"
	ussh "github.com/Hellcatlk/networkconfiguration-operator/pkg/utils/ssh"
)

// NewSSH return openvswitch ssh backend
func NewSSH(ctx context.Context, Host string, cert *certificate.Certificate, options map[string]string) (sw devices.Switch, err error) {
	if len(options) == 0 || options["bridge"] == "" {
		return nil, fmt.Errorf("bridge of openvswitch cli backend is required")
	}

	return &ssh{
		Host:     Host,
		username: cert.Username,
		password: cert.Password,
		bridge:   options["bridge"],
	}, nil
}

// SSH control openvswitch by ssh and cli
type ssh struct {
	Host     string
	username string
	password string
	bridge   string
}

// PowerOn enable openvswitch
func (c *ssh) PowerOn(ctx context.Context) (err error) {
	err = ussh.Run(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "list", "br", c.bridge,
	)) // #nosec
	if err != nil {
		return fmt.Errorf("check birdge failed: %s", err)
	}

	return nil
}

// PowerOff disable openvswitch
func (c *ssh) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr get the port's configure
func (c *ssh) GetPortAttr(ctx context.Context, name string) (configuration *v1alpha1.SwitchPortConfiguration, err error) {
	output, err := ussh.Output(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "list", "port", name,
		"|", "grep", "-E", "-w", "^tag",
		"|", "grep", "-o", "[0-9]*",
	)) // #nosec
	if err != nil {
		return nil, fmt.Errorf("get port failed: %s[%s]", output, err)
	}

	id, err := strconv.Atoi(strings.Trim(string(output), "\n"))
	if err != nil {
		return &v1alpha1.SwitchPortConfiguration{}, nil
	}

	return &v1alpha1.SwitchPortConfiguration{
		Spec: v1alpha1.SwitchPortConfigurationSpec{
			UntaggedVLAN: &id,
		},
	}, nil
}

// SetPortAttr set configure to the port
func (c *ssh) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	if configuration.Spec.UntaggedVLAN == nil {
		return nil
	}

	output, err := ussh.Output(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "set", "port", name, "tag="+strconv.Itoa(*configuration.Spec.UntaggedVLAN),
	)) // #nosec
	if err != nil {
		return fmt.Errorf("set port failed: %s[%s]", output, err)
	}

	actualConfiguration, err := c.GetPortAttr(ctx, name)
	if err != nil {
		return fmt.Errorf("get port failed: %s", err)
	}

	if !reflect.DeepEqual(configuration, actualConfiguration) {
		return fmt.Errorf("set port failed: the actual configuration is inconsistent with the target configuration")
	}

	return nil
}

// ResetPort remove all configure of the port
func (c *ssh) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) (err error) {
	if configuration == nil {
		return nil
	}

	if configuration.Spec.UntaggedVLAN == nil {
		return nil
	}

	output, err := ussh.Output(c.Host, c.username, c.password, exec.Command(
		"ovs-vsctl", "remove ", "port", name, "tag", strconv.Itoa(*configuration.Spec.UntaggedVLAN),
	)) // #nosec
	if err != nil {
		return fmt.Errorf("set port failed: %s[%s]", output, err)
	}

	return nil
}
