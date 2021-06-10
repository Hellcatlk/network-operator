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

	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/certificate"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OVSSwitchSpec defines the desired state of OVSSwitch
type OVSSwitchSpec struct {
	Host   string `json:"host"`
	Bridge string `json:"bridge"`

	// The secret containing the switch credentials
	Secret *corev1.SecretReference `json:"secret"`
}

// OVSSwitchStatus defines the observed state of OVSSwitch
type OVSSwitchStatus struct {
}

// +kubebuilder:object:root=true

// OVSSwitch is the Schema for the ovsswitches API
type OVSSwitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OVSSwitchSpec   `json:"spec,omitempty"`
	Status OVSSwitchStatus `json:"status,omitempty"`
}

// GetConfiguration generate configuration from openvswitch switch
func (s *OVSSwitch) GetConfiguration(ctx context.Context, client client.Client) (*provider.Config, error) {
	cert, err := certificate.Fetch(ctx, client, s.Spec.Secret)
	if err != nil {
		return nil, err
	}

	return &provider.Config{
		OS:       "openvswitch",
		Protocol: "ssh",
		Host:     s.Spec.Host,
		Cert:     cert,
		Options: map[string]string{
			"bridge": s.Spec.Bridge,
		},
	}, nil
}

// +kubebuilder:object:root=true

// OVSSwitchList contains a list of OVSSwitch
type OVSSwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OVSSwitch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OVSSwitch{}, &OVSSwitchList{})
}
