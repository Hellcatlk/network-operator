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

	"github.com/Hellcatlk/networkconfiguration-operator/pkg/provider"
	corev1 "k8s.io/api/core/v1"
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

type testSwitch struct {
}

// GetOS  return switch's os
func (s *testSwitch) GetOS() string {
	return "test"
}

// GetProtocol return switch's protocol
func (s *testSwitch) GetProtocol() string {
	return "test"
}

// GetHost return switch's host
func (s *testSwitch) GetHost() string {
	return ""
}

// GetSecret return switch's certificate secret reference
func (s *testSwitch) GetSecret() *corev1.SecretReference {
	return &corev1.SecretReference{}
}

// GetOptions return switch's options
func (s *testSwitch) GetOptions() map[string]string {
	return nil
}

// Fetch the instance
func (ref *ProviderSwitchRef) Fetch(ctx context.Context, client client.Client) (instance provider.Switch, err error) {
	if ref == nil {
		return nil, fmt.Errorf("provider switch reference is nil")
	}

	switch ref.Kind {
	case "TestSwitch":
		instance = &testSwitch{}
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
	// The reference of provider switch
	ProviderSwitch *ProviderSwitchRef `json:"providerSwitch,omitempty"`

	// Restricted ports in the switch
	Ports map[string]Port `json:"ports,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

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
