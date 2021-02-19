/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StateType is the type of .status.state
type StateType string

// NetworkConfiguration specify the network configuration that metal3Machine needs to use
type NetworkConfiguration struct {
	ConfigurationRefs []ConfigurationRef `json:"configurationRefs"`
	NicHint           NicHint            `json:"nicHint"`
}

// NicHint describes the requirements for the network card
type NicHint struct {
	// The name of the network card for this NicHint.
	Name string `json:"name"`

	// True if smart network card is required, false otherwise.
	SmartNic bool `json:"smartNic"`
}

// PortRef is the reference for Port CR
type PortRef struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	NameSpace string `json:"namespace,omitempty"`

	APIVersions string `json:"apiVersions"`
}

// Fetch the instance
func (ref *PortRef) Fetch(ctx context.Context, client client.Client) (instance *Port, err error) {
	err = client.Get(
		ctx,
		types.NamespacedName{
			Name:      ref.Name,
			Namespace: ref.NameSpace,
		},
		instance,
	)

	return
}

// PortSpec defines the desired state of Port
type PortSpec struct {
	// Reference for PortConfiguration CR.
	ConfigurationRef ConfigurationRef `json:"portConfigurationRef"`

	// Describes the port number on the device.
	ID string `json:"id"`

	// Reference for Port CR.
	// Represents the next port information of this port link
	NextRef PortRef `json:"nextRef,omitempty"`
}

// ConfigurationRef is the reference for Configuration CR
type ConfigurationRef struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	NameSpace string `json:"namespace,omitempty"`

	// +kubebuilder:validation:Enum="SwitchPortConfiguration"
	Kind string `json:"kind"`
}

// Fetch the instance
func (ref *ConfigurationRef) Fetch(ctx context.Context, client client.Client) (instance interface{}, err error) {
	switch ref.Kind {
	case "SwitchPortConfiguration":
		switchPortConfiguration := &SwitchPortConfiguration{}
		err = client.Get(
			ctx,
			types.NamespacedName{
				Name:      ref.Name,
				Namespace: ref.NameSpace,
			},
			switchPortConfiguration,
		)
		instance = switchPortConfiguration
	default:
		err = fmt.Errorf("no instance for the ref")
	}

	return
}

// DeviceRef is the reference for Device CR
type DeviceRef struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	NameSpace string `json:"namespace,omitempty"`

	// +kubebuilder:validation:Enum="Switch"
	Kind string `json:"kind"`
}

// Fetch the instance
func (ref *DeviceRef) Fetch(ctx context.Context, client client.Client) (instance interface{}, err error) {
	switch ref.Kind {
	case "Switch":
		err = client.Get(
			ctx,
			types.NamespacedName{
				Name:      ref.Name,
				Namespace: ref.NameSpace,
			},
			instance.(*Switch),
		)
	default:
		err = fmt.Errorf("no instance for the ref")
	}

	return
}

// PortStatus defines the observed state of Port
type PortStatus struct {
	// The current configuration status of the port.
	State StateType `json:"state,omitempty"`

	// The current Configuration of the port.
	ConfigurationRef ConfigurationRef `json:"ConfigurationRef"`
}

const (
	// PortNone means the port can be configured.
	PortNone StateType = ""

	// PortCreated means we are configuring configuration for the port.
	PortCreated StateType = "Created"

	// PortConfiguring means we are removing configuration from the port.
	PortConfiguring StateType = "Configuring"

	// PortConfigured means the port have been configured, you can use it now.
	PortConfigured StateType = "Configured"

	// PortCleaning means the port have been configured, you can use it now.
	PortCleaning StateType = "Cleaning"

	// PortCleaned means now configuration of the port have been removed.
	PortCleaned StateType = "Cleaned"
)

// GetState gets the current state of the port
func (n *Port) GetState() StateType {
	return n.Status.State
}

// SetState sets the state of the port
func (n *Port) SetState(state StateType) {
	n.Status.State = state
}

// +kubebuilder:object:root=true

// Port is the Schema for the ports API
type Port struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PortSpec   `json:"spec,omitempty"`
	Status PortStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PortList contains a list of Port
type PortList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Port `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Port{}, &PortList{})
}
