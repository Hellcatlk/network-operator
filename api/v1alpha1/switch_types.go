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
	"github.com/Hellcatlk/network-operator/pkg/utils/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Port indicates the specific restriction on the port
type Port struct {
	// Describes the port name on the device
	Name string `json:"name"`

	// True if this port is not available, false otherwise
	Disabled bool `json:"disabled,omitempty"`

	// True if this port can be used as a trunk port, false otherwise
	TrunkDisabled bool `json:"trunkDisable,omitempty"`

	// Indicates the range of VLANs allowed by this port in the switch
	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	// +kubebuilder:default:="1-4096"
	VLANRange string `json:"vlanRange,omitempty"`
}

// Verify configuration
func (p *Port) Verify(configuration *SwitchPortConfiguration) error {
	if p == nil {
		return fmt.Errorf("the port is nil")
	}

	if p.Name == "" {
		return fmt.Errorf("the port's name can't be empty")
	}

	if configuration == nil {
		return nil
	}

	if p.Disabled {
		return fmt.Errorf("the port is disabled")
	}

	if p.TrunkDisabled && configuration.Spec.VLANs != "" {
		return fmt.Errorf("the port can be used as a trunk port")
	}

	if p.VLANRange != "" {
		// Get allowed vlan range
		vlanRange, err := strings.RangeToSlice(p.VLANRange)
		if err != nil {
			return err
		}
		allowed := make(map[int]struct{})
		for _, vlan := range vlanRange {
			allowed[vlan] = struct{}{}
		}

		// Get target vlan range
		target, err := strings.RangeToSlice(p.VLANRange)
		if err != nil {
			return err
		}
		if configuration.Spec.UntaggedVLAN != nil {
			target = append(target, *configuration.Spec.UntaggedVLAN)
		}

		// Check vlan range
		for _, vlan := range target {
			_, existed := allowed[vlan]
			if !existed {
				return fmt.Errorf("vlan %d is out of permissible range", vlan)
			}
		}
	}

	return nil
}

// SwitchProviderRef is the reference for SwitchProvider CR
type SwitchProviderRef struct {
	// +kubebuilder:validation:Enum=AnsibleSwitch
	Kind string `json:"kind"`

	Name string `json:"name"`

	// If empty use default namespace.
	// +kubebuilder:default:="default"
	Namespace string `json:"namespace,omitempty"`
}

// Fetch the instance
func (ref *SwitchProviderRef) Fetch(ctx context.Context, client client.Client) (provider.Switch, error) {
	if ref == nil {
		return nil, fmt.Errorf("provider reference is nil")
	}

	var instance provider.Switch
	var err error

	switch ref.Kind {
	case "TestSwitch":
		instance = &provider.TestSwitch{}

	case "AnsibleSwitch":
		a := &AnsibleSwitch{}
		err = client.Get(
			ctx,
			types.NamespacedName{
				Name:      ref.Name,
				Namespace: ref.Namespace,
			},
			a,
		)
		instance = a

	default:
		err = fmt.Errorf("unknown provider switch kind")
	}

	return instance, err
}

// SwitchSpec defines the desired state of Switch
type SwitchSpec struct {
	// The reference of provider
	Provider *SwitchProviderRef `json:"provider"`

	// Restricted ports in the switch
	Ports map[string]*Port `json:"ports,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
	// The current configuration status of the switch
	State machine.StateType `json:"state,omitempty"`

	// The reference of switch provider
	Provider *SwitchProviderRef `json:"provider,omitempty"`

	// Restricted ports in the switch
	Ports map[string]*Port `json:"ports,omitempty"`

	// The error message of the port
	Error string `json:"error,omitempty"`
}

const (
	// SwitchNone means the CR has just been created
	SwitchNone machine.StateType = ""

	// SwitchVerifying means we are verifying the connection of switch
	SwitchVerifying machine.StateType = "Verifying"

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
	if err != nil {
		s.Status.Error = err.Error()
		return
	}
	s.Status.Error = ""
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
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
