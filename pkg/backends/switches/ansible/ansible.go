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
	"github.com/Hellcatlk/network-operator/pkg/utils/strings"
)

// New return ansible backend
func New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	if config == nil {
		return nil, fmt.Errorf("configure of switch(%s) is nil", config.OS)
	}

	if config.Cert == nil {
		return nil, fmt.Errorf("certificate of switch(%s) is nil", config.OS)
	}

	return &ansible{
		host:   config.Host,
		cert:   config.Cert,
		os:     config.OS,
		bridge: config.Options["bridge"].(string),
	}, nil
}

// ansible backend
type ansible struct {
	host   string
	os     string
	cert   *certificate.Certificate
	bridge string
}

type networkRunnerData struct {
	Host string                   `json:"host"`
	Cert *certificate.Certificate `json:"cert"`
	OS   string                   `json:"os"`
	// Bridge just use for openvswitch
	Bridge       string `json:"bridge,omitempty"`
	Operator     string `json:"operator"`
	Port         string `json:"port"`
	UntaggedVLAN *int   `json:"untaggedVLAN,omitempty"`
	VLANs        []int  `json:"vlans,omitempty"`
}

func (a *ansible) configureAccessPort(port string, untaggedVLAN *int) error {
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

func (a *ansible) configureTrunkPort(port string, untaggedVLAN *int, vlans []int) error {
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

func (a *ansible) deletePort(port string) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Bridge:   a.bridge,
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

// GetPortAttr return the port's configuration
func (a *ansible) GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfiguration, error) {
	return nil, fmt.Errorf("ansible backend does not support GetPortAttr")
}

// SetPortAttr set the configuration to the port
func (a *ansible) SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	if len(configuration.Spec.VLANs) == 0 {
		return a.configureAccessPort(port, configuration.Spec.UntaggedVLAN)
	}

	vlans, err := strings.RangeToSlice(configuration.Spec.VLANs)
	if err != nil {
		return err
	}
	return a.configureTrunkPort(port, configuration.Spec.UntaggedVLAN, vlans)
}

// ResetPort clean the configuration in the port
func (a *ansible) ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return a.deletePort(port)
}
