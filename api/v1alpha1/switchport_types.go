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

	"github.com/Hellcatlk/networkconfiguration-operator/pkg/machine"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SwitchPortReference is the reference for SwitchPort CR
type SwitchPortReference struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	Namespace string `json:"namespace,omitempty"`
}

// Fetch the instance
func (ref *SwitchPortReference) Fetch(ctx context.Context, client client.Client) (instance *SwitchPort, err error) {
	if ref == nil {
		return nil, fmt.Errorf("switch port configuration reference is nil")
	}

	instance = &SwitchPort{}
	err = client.Get(
		ctx,
		types.NamespacedName{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
		instance,
	)

	return
}

// SwitchPortConfigurationReference is the reference for SwitchPortConfiguration CR
type SwitchPortConfigurationReference struct {
	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	Namespace string `json:"namespace,omitempty"`
}

// Fetch the instance
func (ref *SwitchPortConfigurationReference) Fetch(ctx context.Context, client client.Client) (instance *SwitchPortConfiguration, err error) {
	if ref == nil {
		return nil, fmt.Errorf("switch port configuration reference is nil")
	}

	instance = &SwitchPortConfiguration{}
	err = client.Get(
		ctx,
		types.NamespacedName{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
		instance,
	)

	return
}

// SwitchPortSpec defines the desired state of SwitchPort
type SwitchPortSpec struct {
	// Reference for PortConfiguration CR.
	ConfigurationRef *SwitchPortConfigurationReference `json:"configurationRef,omitempty"`
}

// SwitchPortStatus defines the observed state of SwitchPort
type SwitchPortStatus struct {
	// The current configuration status of the port.
	State machine.StateType `json:"state,omitempty"`

	// The error message of the port
	Error string `json:"error,omitempty"`

	// The current Configuration of the port.
	Configuration *SwitchPortConfiguration `json:"configuration,omitempty"`
}

const (
	// SwitchPortNone means the port can be configured
	SwitchPortNone machine.StateType = ""

	// SwitchPortIdle means we are wait for configuration for the port
	SwitchPortIdle machine.StateType = "Idle"

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
	if sp == nil || len(sp.OwnerReferences) == 0 {
		return nil, fmt.Errorf("switch port reference is nil")
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
// +kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".status.state",description="state"
// +kubebuilder:printcolumn:name="ERROR",type="string",JSONPath=".status.error",description="error"

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
