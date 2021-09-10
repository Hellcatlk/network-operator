package ansible

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	ustrings "github.com/Hellcatlk/network-operator/pkg/utils/strings"
	"golang.org/x/crypto/ssh"
)

// New return ansible backend
func New(ctx context.Context, config *provider.SwitchConfiguration) (backends.Switch, error) {
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

type portConfiguration struct {
	Mode         string `json:"mode"`
	VLAN         *int   `json:"vlan,omitempty"`
	TrunkedVLANs string `json:"trunked_vlans,omitempty"`
}

func (a *ansible) getPortConf(port string) (*portConfiguration, error) {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Bridge:   a.bridge,
		Operator: "GetPortConf",
		Port:     port,
	})
	if err != nil {
		return nil, err
	}

	// Execute network runner
	cmd := exec.Command("network-runner", string(data), "2>&1") // #nosec
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s[%s]", output, err)
	}

	// Find last json string from output
	output, err = ustrings.LastJSON(string(output))
	if err != nil {
		return nil, err
	}
	portConfiguration := &portConfiguration{}
	err = json.Unmarshal(output, portConfiguration)
	if err != nil {
		return nil, err
	}

	return portConfiguration, nil
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

	cmd := exec.Command("network-runner", string(data), "2>&1") // #nosec
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

	cmd := exec.Command("network-runner", string(data), "2>&1") // #nosec
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

	cmd := exec.Command("network-runner", string(data), "2&>1") // #nosec
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s[%s]", output, err)
	}

	return nil
}

// IsAvaliable check switch is avaliable or not
func (a *ansible) IsAvaliable() error {
	config := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.Password(a.cert.Password),
		},
		User: a.cert.Username,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	config.SetDefaults()

	address := a.host
	if !strings.Contains(address, ":") {
		address = a.host + ":22"
	}
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return nil
}

// GetPortAttr return the port's configuration
func (a *ansible) GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfigurationSpec, error) {
	portConfiguration, err := a.getPortConf(port)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.SwitchPortConfigurationSpec{
		UntaggedVLAN: portConfiguration.VLAN,
		VLANs:        portConfiguration.TrunkedVLANs,
	}, nil
}

// SetPortAttr set the configuration to the port
func (a *ansible) SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfigurationSpec) error {
	if configuration.VLANs == "" {
		return a.configureAccessPort(port, configuration.UntaggedVLAN)
	}

	vlans, err := ustrings.RangeToSlice(configuration.VLANs)
	if err != nil {
		return err
	}
	return a.configureTrunkPort(port, configuration.UntaggedVLAN, vlans)
}

// ResetPort clean the configuration in the port
func (a *ansible) ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfigurationSpec) error {
	return a.deletePort(port)
}
