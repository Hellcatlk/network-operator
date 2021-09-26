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

	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AnsibleSwitchSpec defines the desired state of AnsibleSwitch
type AnsibleSwitchSpec struct {
	// +kubebuilder:validation:Enum=openvswitch;junos;nxos;eos;enos;cumulus;dellos10;fos
	OS string `json:"os"`

	Host string `json:"host"`

	// A secret containing the switch credentials
	// The default namespace is the same as `AnsibleSwitch`
	Credentials *corev1.SecretReference `json:"credentials"`

	// OVS bridge
	Bridge string `json:"bridge,omitempty"`
}

// AnsibleSwitchStatus defines the observed state of AnsibleSwitch
type AnsibleSwitchStatus struct {
}

// +kubebuilder:object:root=true

// AnsibleSwitch is the Schema for the ansibleswitches API
type AnsibleSwitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnsibleSwitchSpec   `json:"spec,omitempty"`
	Status AnsibleSwitchStatus `json:"status,omitempty"`
}

// GetConfiguration generate configuration from openvswitch switch
func (a *AnsibleSwitch) GetConfiguration(ctx context.Context, client client.Client) (*provider.SwitchConfiguration, error) {
	// Set the default namespace of `Credentials` to the same as `AnsibleSwitch`
	if a.Spec.Credentials.Namespace == "" {
		a.Spec.Credentials.Namespace = a.Namespace
	}

	cert, err := credentials.Fetch(ctx, client, a.Spec.Credentials)
	if err != nil {
		return nil, err
	}

	if a.Spec.OS == "openvswitch" && a.Spec.Bridge == "" {
		return nil, fmt.Errorf("for openvswitch bridge is required")
	}

	return &provider.SwitchConfiguration{
		OS:          a.Spec.OS,
		Host:        a.Spec.Host,
		Credentials: cert,
		Backend:     "ansible",
		Options: map[string]interface{}{
			"bridge": a.Spec.Bridge,
		},
	}, nil
}

// +kubebuilder:object:root=true

// AnsibleSwitchList contains a list of AnsibleSwitch
type AnsibleSwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnsibleSwitch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnsibleSwitch{}, &AnsibleSwitchList{})
}
