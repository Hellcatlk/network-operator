package ssh

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	ussh "github.com/Hellcatlk/network-operator/pkg/utils/ssh"
)

// SSH control openvswitch by ssh and cli
type SSH struct {
	Host   string
	cert   *certificate.Certificate
	bridge string
}

// New return ssh backend
func (s *SSH) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	if config == nil {
		return nil, fmt.Errorf("configure of switch(%s) is nil", config.OS)
	}

	if config.OS != "openvswitch" {
		return nil, fmt.Errorf("currently the ssh backend only supports openvswitch")
	}

	if config.Cert == nil {
		return nil, fmt.Errorf("certificate of switch(%s) is nil", config.OS)
	}

	return &SSH{
		Host:   config.Host,
		cert:   config.Cert,
		bridge: config.Options["bridge"],
	}, nil
}

// PowerOn enable openvswitch
func (s *SSH) PowerOn(ctx context.Context) error {
	err := ussh.Run(s.Host, s.cert.Username, s.cert.Password, exec.Command(
		"sudo", "ovs-vsctl", "list", "br", s.bridge,
	)) // #nosec
	if err != nil {
		return fmt.Errorf("check birdge failed: %s", err)
	}

	return nil
}

// PowerOff disable openvswitch
func (s *SSH) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr get the port's configure
func (s *SSH) GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfiguration, error) {
	output, err := ussh.Output(s.Host, s.cert.Username, s.cert.Password, exec.Command(
		"sudo", "ovs-vsctl", "list", "port", port,
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
func (s *SSH) SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	if configuration == nil {
		return nil
	}

	if configuration.Spec.UntaggedVLAN == nil {
		return nil
	}

	output, err := ussh.Output(s.Host, s.cert.Username, s.cert.Password, exec.Command(
		"sudo", "ovs-vsctl", "set", "port", port, "tag="+strconv.Itoa(*configuration.Spec.UntaggedVLAN),
	)) // #nosec
	if err != nil {
		return fmt.Errorf("set port failed: %s[%s]", output, err)
	}

	actualConfiguration, err := s.GetPortAttr(ctx, port)
	if err != nil {
		return fmt.Errorf("get port failed: %s", err)
	}

	if !reflect.DeepEqual(configuration.Spec, actualConfiguration.Spec) {
		return fmt.Errorf("set port failed: the actual configuration is inconsistent with the target configuration")
	}

	return nil
}

// ResetPort remove all configure of the port
func (s *SSH) ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	if configuration == nil {
		return nil
	}

	if configuration.Spec.UntaggedVLAN == nil {
		return nil
	}

	output, err := ussh.Output(s.Host, s.cert.Username, s.cert.Password, exec.Command(
		"sudo", "ovs-vsctl", "remove ", "port", port, "tag", strconv.Itoa(*configuration.Spec.UntaggedVLAN),
	)) // #nosec
	if err != nil {
		return fmt.Errorf("set port failed: %s[%s]", output, err)
	}

	return nil
}
