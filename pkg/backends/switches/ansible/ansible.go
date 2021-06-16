package ansible

import (
	"context"
	"fmt"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/sbinet/go-python"
)

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

// Ansible backend
type Ansible struct {
	networkRunner *python.PyObject
}

func (a *Ansible) createVlan(vlan int) error {
	a.networkRunner.CallMethod("create_vlan", "ansibleHost", vlan)
	return nil
}

func (a *Ansible) deleteVlan(vlan int) error {
	a.networkRunner.CallMethod("delete_vlan", "ansibleHost", vlan)
	return nil
}

func (a *Ansible) configureAccessPort(port string, vlan int) error {
	a.networkRunner.CallMethod("conf_access_port", "ansibleHost", "port", vlan)
	return nil
}

func (a *Ansible) configureTrunkPort(port string, vlans []int) error {
	a.networkRunner.CallMethod("conf_trunk_port", "ansibleHost", "port", vlans)
	return nil
}

func (a *Ansible) deletePort(port string) error {
	a.networkRunner.CallMethod("delete_port", "ansibleHost", "port")
	return nil
}

// New return test backend
func (a *Ansible) New(ctx context.Context, config *provider.Config) (backends.Switch, error) {
	// Import modules
	ModuleInventory := python.PyImport_ImportModule("network_runner.models.inventory")
	if ModuleInventory == nil {
		return nil, fmt.Errorf("import python module network_runner.models.inventory failed")
	}
	ModuleAPI := python.PyImport_ImportModule("network_runner.api")
	if ModuleAPI == nil {
		return nil, fmt.Errorf("import python module network_runner.api failed")
	}

	// Get python classes
	classHost := ModuleInventory.GetAttrString("Host")
	if classHost == nil {
		return nil, fmt.Errorf("get python class Host failed")
	}
	classInventory := ModuleInventory.GetAttrString("Inventory")
	if classInventory == nil {
		return nil, fmt.Errorf("get python class Inventory failed")
	}
	classNetworkRunner := ModuleAPI.GetAttrString("NetworkRunner")
	if classNetworkRunner == nil {
		return nil, fmt.Errorf("get python class NetworkRunner failed")
	}

	// Network runner initial
	host := python.PyInstance_New(classHost, nil, nil)
	if host == nil {
		return nil, fmt.Errorf("initital python object Host failed")
	}
	inventory := python.PyInstance_New(classInventory, nil, nil)
	if inventory == nil {
		return nil, fmt.Errorf("initital python object Inventory failed")
	}
	inventory.GetAttrString("hosts").CallMethodObjArgs("add", host)
	networkRunner := python.PyInstance_New(classNetworkRunner, inventory, nil)
	if networkRunner == nil {
		return nil, fmt.Errorf("initital python object NetworkRunner failed")
	}

	return &Ansible{
		networkRunner: networkRunner,
	}, nil
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
func (a *Ansible) GetPortAttr(ctx context.Context, name string) (*v1alpha1.SwitchPortConfiguration, error) {
	return &v1alpha1.SwitchPortConfiguration{}, nil
}

// SetPortAttr just for test
func (a *Ansible) SetPortAttr(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}

// ResetPort just for test
func (a *Ansible) ResetPort(ctx context.Context, name string, configuration *v1alpha1.SwitchPortConfiguration) error {
	return nil
}
