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
	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AnsibleSpec defines the desired state of Ansible
type AnsibleSpec struct {
	// +kubebuilder:validation:Enum=openvswitch;junos;nxos;eos;enos;cumulus;dellos10;fos
	OS string `json:"os"`

	Host string `json:"host"`

	// OVS bridge
	Bridge string `json:"bridge,omitempty"`

	// The secret containing the switch credentials
	Secret *corev1.SecretReference `json:"secret"`
}

// AnsibleStatus defines the observed state of Ansible
type AnsibleStatus struct {
}

// +kubebuilder:object:root=true

// Ansible is the Schema for the ansibles API
type Ansible struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnsibleSpec   `json:"spec,omitempty"`
	Status AnsibleStatus `json:"status,omitempty"`
}

// GetConfiguration generate configuration from openvswitch switch
func (a *Ansible) GetConfiguration(ctx context.Context, client client.Client) (*provider.Config, error) {
	cert, err := certificate.Fetch(ctx, client, a.Spec.Secret)
	if err != nil {
		return nil, err
	}

	if a.Spec.OS == "openvswitch" && a.Spec.Bridge == "" {
		return nil, fmt.Errorf("for openvswitch bridge is required")
	}

	return &provider.Config{
		OS:      a.Spec.OS,
		Host:    a.Spec.Host,
		Cert:    cert,
		Backend: "ansible",
		Options: map[string]interface{}{
			"bridge": a.Spec.Bridge,
		},
	}, nil
}

// +kubebuilder:object:root=true

// AnsibleList contains a list of Ansible
type AnsibleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ansible `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ansible{}, &AnsibleList{})
}
