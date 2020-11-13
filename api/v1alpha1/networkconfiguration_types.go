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

func init() {
	SchemeBuilder.Register(&NetworkConfiguration{}, &NetworkConfigurationList{})
}

// +kubebuilder:object:root=true

// NetworkConfigurationList contains a list of NetworkConfiguration
type NetworkConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkConfiguration `json:"items"`
}

// +kubebuilder:object:root=true

// NetworkConfiguration is the Schema for the networkconfigurations API
type NetworkConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkConfigurationSpec   `json:"spec,omitempty"`
	Status NetworkConfigurationStatus `json:"status,omitempty"`
}

// NetworkConfigurationSpec defines the desired state of NetworkConfiguration
type NetworkConfigurationSpec struct {
	// +kubebuilder:validation:MaxItems=10
	ACLs []ACL `json:"acls,omitempty"`

	// +kubebuilder:validation:MaxItems=3
	Vlans []string `json:"vlans,omitempty"`

	UntaggedVLAN string `json:"untaggedVLAN,omitempty"`

	// +kubebuilder:validation:Enum="lag";"mlag"
	LinkAggregationType string `json:"linkAggregationType,omitempty"`

	NICHint NICHint `json:"nicHint,omitempty"`
}

// ACL ...
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

// NetworkConfigurationStatus defines the observed state of NetworkConfiguration
type NetworkConfigurationStatus struct {
	NetworkBindingRefs []NetworkBindingRef `json:"networkBindingRefs,omitempty"`
}
