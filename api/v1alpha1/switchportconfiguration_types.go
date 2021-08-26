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
	"reflect"

	"github.com/Hellcatlk/network-operator/pkg/utils/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ACL describes the rules applied in the switch
type ACL struct {
	// +kubebuilder:validation:Enum=ipv4;ipv6
	Type string `json:"type,omitempty"`

	// +kubebuilder:validation:Enum=allow;deny
	Action string `json:"action,omitempty"`

	// +kubebuilder:validation:Enum=TCP;UDP;ICMP;ALL
	Protocol string `json:"protocol,omitempty"`

	Src string `json:"src,omitempty"`

	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	SrcPortRange string `json:"srcPortRange,omitempty"`

	Des string `json:"des,omitempty"`

	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	DesPortRange string `json:"desPortRange,omitempty"`
}

// SwitchPortConfigurationSpec defines the desired state of SwitchPortConfiguration
type SwitchPortConfigurationSpec struct {
	// +kubebuilder:validation:MaxItems=10
	ACLs []ACL `json:"acls,omitempty"`

	UntaggedVLAN *int `json:"untaggedVLAN,omitempty"`

	// The range of tagged vlans. You can use `-` to connect two numbers to express the range
	// or use separate numbers. You can use `,` to combine the above two methods, for example:
	// `1-10,11,13-20`
	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	VLANs string `json:"vlans,omitempty"`

	// Disable port
	Disable bool `json:"disable,omitempty"`
}

// IsEqual check configuration is equal or not
func (target *SwitchPortConfigurationSpec) IsEqual(actual *SwitchPortConfigurationSpec) bool {
	if target == actual {
		return true
	}
	if target == nil && reflect.DeepEqual(actual, &SwitchPortConfigurationSpec{}) {
		return true
	}
	if actual == nil && reflect.DeepEqual(target, &SwitchPortConfigurationSpec{}) {
		return true
	}

	rangeTarget, err := strings.RangeToSlice(target.VLANs)
	if err != nil {
		return false
	}
	rangeActual, err := strings.RangeToSlice(actual.VLANs)
	if err != nil {
		return false
	}
	if !reflect.DeepEqual(rangeTarget, rangeActual) {
		return false
	}

	targetCopy := target.DeepCopy()
	targetCopy.VLANs = ""
	actualCopy := actual.DeepCopy()
	actualCopy.VLANs = ""
	return reflect.DeepEqual(targetCopy, actualCopy)
}

// SwitchPortConfigurationStatus defines the observed state of SwitchPortConfiguration
type SwitchPortConfigurationStatus struct {
}

// +kubebuilder:object:root=true

// SwitchPortConfiguration is the Schema for the switchportconfigurations API
type SwitchPortConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchPortConfigurationSpec   `json:"spec,omitempty"`
	Status SwitchPortConfigurationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SwitchPortConfigurationList contains a list of SwitchPortConfiguration
type SwitchPortConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchPortConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchPortConfiguration{}, &SwitchPortConfigurationList{})
}
