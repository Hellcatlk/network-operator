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
	"github.com/Hellcatlk/network-operator/pkg/utils/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Verify that the allocated resources are within the allowed range
func (l *TenantLimit) Verify(available *SwitchResourceStatus) error {
	if l == nil {
		return nil
	}

	if l.VLANRange == "" {
		return nil
	}

	// Get allowed vlan range
	vlanRange, err := strings.RangeToSlice(available.AvailableVLAN)
	if err != nil {
		return err
	}
	allowed := make(map[int]struct{})
	for _, vlan := range vlanRange {
		allowed[vlan] = struct{}{}
	}

	// Get target vlan range
	target, err := strings.RangeToSlice(l.VLANRange)
	if err != nil {
		return err
	}

	// Check vlan range
	for _, vlan := range target {
		_, existed := allowed[vlan]
		if !existed {
			return fmt.Errorf("vlan %d is out of permissible range", vlan)
		}
	}

	return nil
}

// Expansion ...
func (sr *SwitchResource) Expansion(limit *TenantLimit) error {
	if limit == nil {
		return nil
	}

	result, err := strings.Expansion(sr.Status.AvailableVLAN, limit.VLANRange)
	if err != nil {
		return err
	}

	sr.Status.AvailableVLAN = result

	return nil
}

// Shrink ...
func (sr *SwitchResource) Shrink(limit *TenantLimit) error {
	if limit == nil {
		return nil
	}

	result, err := strings.Shrink(sr.Status.AvailableVLAN, limit.VLANRange)
	if err != nil {
		return err
	}

	sr.Status.AvailableVLAN = result

	return nil
}

// FetchSwitchResourceLimit fetch the SwitchResourceLimit/user-limit instance
func (l *TenantLimit) FetchSwitchResourceLimit(ctx context.Context, client client.Client) (*SwitchResourceLimit, error) {
	if l == nil {
		return nil, fmt.Errorf("TenantLimit is nil")
	}

	instance := &SwitchResourceLimit{}
	err := client.Get(
		ctx,
		types.NamespacedName{
			Name:      "user-limit",
			Namespace: l.Namespace,
		},
		instance,
	)

	return instance, err
}

// GetMetadataAndSpec return metadata and spec field
func (sr *SwitchResource) GetMetadataAndSpec() interface{} {
	deepCopy := sr.DeepCopy()
	deepCopy.Status = SwitchResourceStatus{}
	return deepCopy
}

// GetStatus return status field
func (sr *SwitchResource) GetStatus() interface{} {
	return sr.Status.DeepCopy()
}

// GetState gets the current state of the SwitchResource
func (sr *SwitchResource) GetState() machine.StateType {
	return sr.Status.State
}

// SetState sets the state of the SwitchResource
func (sr *SwitchResource) SetState(state machine.StateType) {
	sr.Status.State = state
}

// SetError sets the error of the SwitchResource
func (sr *SwitchResource) SetError(err error) {
	if err != nil {
		sr.Status.Error = err.Error()
		return
	}
	sr.Status.Error = ""
}

// VerifyConfiguration verify that the configuration meets the limit.
func (l *TenantLimit) VerifyConfiguration(configuration *SwitchPortConfiguration) error {
	if l == nil {
		return nil
	}
	if l.VLANRange == "" {
		return nil
	}
	// Get allowed vlan range
	vlanRange, err := strings.RangeToSlice(l.VLANRange)
	if err != nil {
		return err
	}
	allowed := make(map[int]struct{})
	for _, vlan := range vlanRange {
		allowed[vlan] = struct{}{}
	}

	// Get target vlan range
	target, err := strings.RangeToSlice(l.VLANRange)
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

	return nil
}

const (
	// SwitchResourceNone means the CR has just been created
	SwitchResourceNone machine.StateType = ""

	// SwitchResourceVerifying means we are verifying the restrictions on the user
	SwitchResourceVerifying machine.StateType = "Verifying"

	// SwitchResourceCreating means we are creating SwitchResourceLimit
	SwitchResourceCreating machine.StateType = "Creating"

	// SwitchResourceRunning means all of SwitchResourceLimit have been created
	SwitchResourceRunning machine.StateType = "Running"

	// SwitchResourceDeleting means we are deleting SwitchResourceLimit
	SwitchResourceDeleting machine.StateType = "Deleting"
)

// SwitchResourceSpec defines the desired state of SwitchResource
type SwitchResourceSpec struct {
	// Indicates the initial allocatable vlan range
	VLANRange    string                  `json:"vlanRange,omitempty"`
	TenantLimits map[string]*TenantLimit `json:"tenantLimits,omitempty"`
}

// SwitchResourceStatus defines the observed state of SwitchResource
type SwitchResourceStatus struct {
	// Indicates the vlan range that the administrator
	// can assign to the user currently.
	AvailableVLAN string                  `json:"availableVLAN,omitempty"`
	TenantLimits  map[string]*TenantLimit `json:"tenantLimits,omitempty"`
	// The error message of the port
	Error string `json:"error,omitempty"`
	// The current configuration status of the SwitchResource
	State machine.StateType `json:"state,omitempty"`
}

// TenantLimit indicates resource restrictions on tenants
type TenantLimit struct {
	Namespace string `json:"namespace,omitempty"`
	VLANRange string `json:"vlanRange,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SwitchResource is the Schema for the switchresources API
type SwitchResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchResourceSpec   `json:"spec,omitempty"`
	Status SwitchResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchResourceList contains a list of SwitchResource
type SwitchResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchResource{}, &SwitchResourceList{})
}
