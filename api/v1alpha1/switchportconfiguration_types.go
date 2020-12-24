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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwitchPortConfigurationSpec defines the desired state of SwitchPortConfiguration
type SwitchPortConfigurationSpec struct {
	// +kubebuilder:validation:MaxItems=10
	ACLs []ACL `json:"acls,omitempty"`
	// +kubebuilder:validation:MaxItems=3
	Vlans []VLAN `json:"vlans,omitempty"`
	// The untagged VLAN ID
	VLANID VLANID `json:"vlanId,omitempty"`
	// True if it is to be set to trunk port, false otherwise.
	Trunk bool `json:"trunk,omitempty"`
}

// ACL describes the rules applied in the switch
type ACL struct {
	// +kubebuilder:validation:Enum="ipv4";"ipv6"
	Type string `json:"type,omitempty"`

	// +kubebuilder:validation:Enum="allow";"deny"
	Action string `json:"action,omitempty"`

	// +kubebuilder:validation:Enum="TCP";"UDP";"ICMP";"ALL"
	Protocol string `json:"protocol,omitempty"`

	Src string `json:"src,omitempty"`

	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	SrcPortRange string `json:"srcPortRange,omitempty"`

	Des string `json:"des,omitempty"`

	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	DesPortRange string `json:"desPortRange,omitempty"`
}

// SwitchPortConfigurationStatus defines the observed state of SwitchPortConfiguration
type SwitchPortConfigurationStatus struct {
	PortRefs []PortRef `json:"portRefs"`
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
