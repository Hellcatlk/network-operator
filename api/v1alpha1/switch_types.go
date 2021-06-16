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

	"github.com/Hellcatlk/network-operator/pkg/machine"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Port indicates the specific restriction on the port
type Port struct {
	// Describes the port name on the device
	Name string `json:"name,omitempty"`

	// True if this port is not available, false otherwise
	Disabled bool `json:"disabled,omitempty"`

	// True if this port can be used as a trunk port, false otherwise
	TrunkDisabled bool `json:"trunkDisable,omitempty"`

	// Indicates the range of VLANs allowed by this port in the switch
	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	VlanRange string `json:"vlanRange,omitempty"`
}

// ProviderSwitchRef is the reference for ProviderSwitch CR
type ProviderSwitchRef struct {
	Kind string `json:"kind"`

	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	Namespace string `json:"namespace,omitempty"`
}

// Fetch the instance
func (ref *ProviderSwitchRef) Fetch(ctx context.Context, client client.Client) (provider.Switch, error) {
	if ref == nil {
		return nil, fmt.Errorf("provider switch reference is nil")
	}

	var instance provider.Switch
	var err error

	switch ref.Kind {
	case "TestSwitch":
		instance = &provider.Test{}

	case "OVSSwitch":
		ps := &OVSSwitch{}
		err = client.Get(
			ctx,
			types.NamespacedName{
				Name:      ref.Name,
				Namespace: ref.Namespace,
			},
			ps,
		)
		instance = ps

	default:
		err = fmt.Errorf("unknown provider switch kind")
	}

	return instance, err
}

// SwitchSpec defines the desired state of Switch
type SwitchSpec struct {
	// +kubebuilder:validation:Enum=ssh;ansible
	Backend string `json:"backend"`

	// The reference of provider switch
	ProviderSwitch *ProviderSwitchRef `json:"providerSwitch,omitempty"`

	// Restricted ports in the switch
	Ports map[string]Port `json:"ports,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
	// The current configuration status of the switch.
	State machine.StateType `json:"state,omitempty"`

	// The reference of provider switch
	ProviderSwitch *ProviderSwitchRef `json:"providerSwitch,omitempty"`

	// Restricted ports in the switch
	Ports map[string]Port `json:"ports,omitempty"`

	// The error message of the port
	Error string `json:"error,omitempty"`
}

const (
	// SwitchNone means the CR has just been created
	SwitchNone machine.StateType = ""

	// SwitchVerify means we are verifying the connection of switch
	SwitchVerify machine.StateType = "Verifying"

	// SwitchConfiguring means we are creating SwitchPort
	SwitchConfiguring machine.StateType = "Configuring"

	// SwitchRunning means all of SwitchPort have been created
	SwitchRunning machine.StateType = "Running"

	// SwitchDeleting means we are deleting SwitchPort
	SwitchDeleting machine.StateType = "Deleting"
)

// GetState gets the current state of the switch
func (s *Switch) GetState() machine.StateType {
	return s.Status.State
}

// SetState sets the state of the switch
func (s *Switch) SetState(state machine.StateType) {
	s.Status.State = state
}

// SetError sets the error of the switch
func (s *Switch) SetError(err error) {
	if err == nil {
		s.Status.Error = ""
	} else {
		s.Status.Error = err.Error()
	}
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".status.state",description="state"
// +kubebuilder:printcolumn:name="ERROR",type="string",JSONPath=".status.error",description="error"

// Switch is the Schema for the switches API
type Switch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchSpec   `json:"spec,omitempty"`
	Status SwitchStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SwitchList contains a list of Switch
type SwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Switch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Switch{}, &SwitchList{})
}
