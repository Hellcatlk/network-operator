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

// SwitchResourceLimitSpec defines the desired state of SwitchResourceLimit
type SwitchResourceLimitSpec struct {
	// Indicates the range of VLANs allowed by this user
	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	VLANRange string `json:"vlanRange,omitempty"`
}

// SwitchResourceLimitStatus defines the observed state of SwitchResourceLimit
type SwitchResourceLimitStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SwitchResourceLimit is the Schema for the switchresourcelimits API
type SwitchResourceLimit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchResourceLimitSpec   `json:"spec,omitempty"`
	Status SwitchResourceLimitStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchResourceLimitList contains a list of SwitchResourceLimit
type SwitchResourceLimitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchResourceLimit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchResourceLimit{}, &SwitchResourceLimitList{})
}
