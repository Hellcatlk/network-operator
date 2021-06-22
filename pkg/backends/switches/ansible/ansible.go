package ansible

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
)

// Ansible backend
type Ansible struct {
	host    string
	os      string
	cert    *certificate.Certificate
	options map[string]string
}

type networkRunnerData struct {
	Host string                   `json:"host"`
	Cert *certificate.Certificate `json:"cert"`
	OS   string                   `json:"os"`
	// Bridge just use for openvswitch
	Bridge       string          `json:"bridge,omitempty"`
	Operator     string          `json:"operator"`
	Port         string          `json:"port"`
	UntaggedVLAN *v1alpha1.VLAN  `json:"untaggedVLAN,omitempty"`
	VLANs        []v1alpha1.VLAN `json:"vlans,omitempty"`
}

func (a *Ansible) configureAccessPort(port string, untaggedVLAN *v1alpha1.VLAN) error {
	data, err := json.Marshal(networkRunnerData{
		Host:         a.host,
		Cert:         a.cert,
		OS:           a.os,
		Operator:     "ConfigAccessPort",
		Port:         port,
		UntaggedVLAN: untaggedVLAN,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("network-runner", string(data)) // #nosec
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s[%s]", output, err)
	}

	return nil
}

func (a *Ansible) configureTrunkPort(port string, untaggedVLAN *v1alpha1.VLAN, vlans []v1alpha1.VLAN) error {
	data, err := json.Marshal(networkRunnerData{
		Host:         a.host,
		Cert:         a.cert,
		OS:           a.os,
		Operator:     "ConfigTrunkPort",
		Port:         port,
		UntaggedVLAN: untaggedVLAN,
		VLANs:        vlans,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("network-runner", string(data)) // #nosec
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s[%s]", output, err)
	}

	return nil
}

func (a *Ansible) deletePort(port string) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Bridge:   a.options["bridge"],
		Operator: "DeletePort",
		Port:     port,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("network-runner", string(data)) // #nosec
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s[%s]", output, err)
	}

	return nil
}

// New return ansible backend
func (a *Ansible) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	if config == nil {
		return nil, fmt.Errorf("configure of switch(%s) is nil", config.OS)
	}

	if config.Cert == nil {
		return nil, fmt.Errorf("certificate of switch(%s) is nil", config.OS)
	}

	return &Ansible{
		host:    config.Host,
		cert:    config.Cert,
		os:      config.OS,
		options: config.Options,
	}, nil
}

// GetPortAttr return the port's configuration
func (a *Ansible) GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfiguration, error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr set the configuration to the port
func (a *Ansible) SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	if len(configuration.Spec.VLANs) == 0 {
		return a.configureAccessPort(port, configuration.Spec.UntaggedVLAN)
	}

	return a.configureTrunkPort(port, configuration.Spec.UntaggedVLAN, configuration.Spec.VLANs)
}

// ResetPort clean the configuration in the port
func (a *Ansible) ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return a.deletePort(port)
}
