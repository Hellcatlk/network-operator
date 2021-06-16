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
	host string
	os   string
	cert *certificate.Certificate
}

type networkRunnerData struct {
	Host     string
	Cert     *certificate.Certificate
	OS       string
	Operator string
	Port     string
	Vlan     int
	Vlans    []int
}

func (a *Ansible) createVlan(vlan int) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Operator: "CreateVlan",
		Vlan:     vlan,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("./cmd/network_runner/main.py", string(data))
	return cmd.Run()
}

func (a *Ansible) deleteVlan(vlan int) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Operator: "DeleteVlan",
		Vlan:     vlan,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("./cmd/network_runner/main.py", string(data))
	return cmd.Run()
}

func (a *Ansible) configureAccessPort(port string, vlan int) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Operator: "ConfigAccessPort",
		Port:     port,
		Vlan:     vlan,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("./cmd/network_runner/main.py", string(data))
	return cmd.Run()
}

func (a *Ansible) configureTrunkPort(port string, vlans []int) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Operator: "ConfigAccessPort",
		Port:     port,
		Vlans:    vlans,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("./cmd/network_runner/main.py", string(data))
	return cmd.Run()
}

func (a *Ansible) deletePort(port string) error {
	data, err := json.Marshal(networkRunnerData{
		Host:     a.host,
		Cert:     a.cert,
		OS:       a.os,
		Operator: "DeletePort",
		Port:     port,
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("./cmd/network_runner/main.py", string(data))
	return cmd.Run()
}

// New return test backend
func (a *Ansible) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	if config == nil {
		return nil, fmt.Errorf("configure of switch(%s) is nil", config.OS)
	}

	if config.Cert == nil {
		return nil, fmt.Errorf("certificate of switch(%s) is nil", config.OS)
	}

	return &Ansible{
		host: config.Host,
		cert: config.Cert,
		os:   config.OS,
	}, a.deletePort("11")
}

// PowerOn just for test
func (a *Ansible) PowerOn(ctx context.Context) error {
	return nil
}

// PowerOff just for test
func (a *Ansible) PowerOff(ctx context.Context) (err error) {
	return
}

// GetPortAttr just for test
func (a *Ansible) GetPortAttr(ctx context.Context, port string) (*v1alpha1.SwitchPortConfiguration, error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr just for test
func (a *Ansible) SetPortAttr(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}

// ResetPort just for test
func (a *Ansible) ResetPort(ctx context.Context, port string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return a.deletePort(port)
}
