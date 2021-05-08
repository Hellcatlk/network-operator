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

	"github.com/metal3-io/networkconfiguration-operator/pkg/machine"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SwitchPortRef is the reference for Port CR
type SwitchPortRef struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	NameSpace string `json:"namespace,omitempty"`

	APIVersions string `json:"apiVersions"`
}

// Fetch the instance
func (ref *SwitchPortRef) Fetch(ctx context.Context, client client.Client) (instance *SwitchPort, err error) {
	if ref == nil {
		return nil, fmt.Errorf("reference is nil")
	}

	instance = &SwitchPort{}
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

// SwitchPortSpec defines the desired state of SwitchPort
type SwitchPortSpec struct {
	// Reference for PortConfiguration CR.
	ConfigurationRef *SwitchPortConfigurationRef `json:"configurationRef,omitempty"`

	// Describes the port number on the device.
	ID string `json:"id"`
}

// SwitchPortConfigurationRef is the reference for Configuration CR
type SwitchPortConfigurationRef struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	NameSpace string `json:"namespace,omitempty"`
}

// Fetch the instance
func (ref *SwitchPortConfigurationRef) Fetch(ctx context.Context, client client.Client) (instance *SwitchPortConfiguration, err error) {
	if ref == nil {
		return nil, fmt.Errorf("reference is nil")
	}

	instance = &SwitchPortConfiguration{}
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

// SwitchPortStatus defines the observed state of SwitchPort
type SwitchPortStatus struct {
	// The current configuration status of the port.
	State machine.StateType `json:"state,omitempty"`

	// The current Configuration of the port.
	Configuration *SwitchPortConfiguration `json:"Configuration"`
}

const (
	// SwitchPortNone means the port can be configured
	SwitchPortNone machine.StateType = ""

	// SwitchPortIdle means we are wait for configuration for the port
	SwitchPortIdle machine.StateType = "Idle"

	// SwitchPortValidating means we are validating the connection for network device
	SwitchPortValidating machine.StateType = "Validating"

	// SwitchPortConfiguring means we are removing configuration from the port
	SwitchPortConfiguring machine.StateType = "Configuring"

	// SwitchPortActive means the port have been configured, you can use it now
	SwitchPortActive machine.StateType = "Active"

	// SwitchPortCleaning means the port have been configured, you can use it now
	SwitchPortCleaning machine.StateType = "Cleaning"

	// SwitchPortDeleting means we are deleting this CR
	SwitchPortDeleting machine.StateType = "Deleting"
)

// GetState gets the current state of the port
func (sp *SwitchPort) GetState() machine.StateType {
	return sp.Status.State
}

// SetState sets the state of the port
func (sp *SwitchPort) SetState(state machine.StateType) {
	sp.Status.State = state
}

// FetchOwnerReference fetch OwnerReference[0]
func (sp *SwitchPort) FetchOwnerReference(ctx context.Context, client client.Client) (instance *Switch, err error) {
	if sp == nil {
		return nil, fmt.Errorf("reference is nil")
	}

	instance = &Switch{}
	err = client.Get(
		ctx,
		types.NamespacedName{
			Name:      sp.OwnerReferences[0].Name,
			Namespace: sp.Namespace,
		},
		instance,
	)

	return instance, err
}

// +kubebuilder:object:root=true

// SwitchPort is the Schema for the switchports API
type SwitchPort struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchPortSpec   `json:"spec,omitempty"`
	Status SwitchPortStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SwitchPortList contains a list of SwitchPort
type SwitchPortList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchPort `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchPort{}, &SwitchPortList{})
}
