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
	"strconv"

	"github.com/Hellcatlk/network-operator/pkg/utils/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FetchSwitchResource fetch the SwitchResource instance
func (rl *SwitchResourceLimit) FetchSwitchResource(ctx context.Context, client client.Client) (*SwitchResource, error) {
	if rl == nil {
		return nil, fmt.Errorf("SwitchResourceLimit is nil")
	}

	instance := &SwitchResource{}

	err := client.Get(
		ctx,
		types.NamespacedName{
			Name:      rl.Status.SwitchResourceRef.Name,
			Namespace: rl.Status.SwitchResourceRef.Namespace,
		},
		instance,
	)

	return instance, err
}

func (sr *SwitchResourceLimit) Expansion(configuration *SwitchPortConfigurationSpec) error {
	if configuration == nil {
		return nil
	}
	var err error
	useVLAN := configuration.TaggedVLANRange
	if configuration.UntaggedVLAN != nil {
		useVLAN, err = strings.Expansion(configuration.TaggedVLANRange, strconv.Itoa(*configuration.UntaggedVLAN))
		if err != nil {
			return err
		}

	}

	result, err := strings.Expansion(sr.Status.UsedVLAN, useVLAN)
	if err != nil {
		return err
	}

	sr.Status.UsedVLAN = result

	return nil
}

func (sr *SwitchResourceLimit) Shrink(configuration *SwitchPortConfigurationSpec) error {
	if configuration == nil {
		return nil
	}
	var err error
	useVLAN := configuration.TaggedVLANRange
	if configuration.UntaggedVLAN != nil {
		useVLAN, err = strings.Expansion(configuration.TaggedVLANRange, strconv.Itoa(*configuration.UntaggedVLAN))
		if err != nil {
			return err
		}

	}
	result, err := strings.Shrink(sr.Status.UsedVLAN, useVLAN)
	if err != nil {
		return err
	}

	sr.Status.UsedVLAN = result

	return nil
}

// SwitchResourceLimitSpec defines the desired state of SwitchResourceLimit
type SwitchResourceLimitSpec struct {
}

// SwitchResourceLimitStatus defines the observed state of SwitchResourceLimit
type SwitchResourceLimitStatus struct {
	// Indicates the range of VLANs allowed
	// +kubebuilder:validation:Pattern=`([0-9]{1,})|([0-9]{1,}-[0-9]{1,})(,([0-9]{1,})|([0-9]{1,}-[0-9]{1,}))*`
	// +kubebuilder:default:="1-4096"
	VLANRange         string            `json:"vlanRange,omitempty"`
	SwitchResourceRef SwitchResourceRef `json:"switchResourceRef,omitempty"`
	UsedVLAN          string            `json:"usedVLAN,omitempty"`
}

type SwitchResourceRef struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
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
