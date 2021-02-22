package controllers

import (
	"context"
	"testing"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/machine"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type fakeClient struct {
	client.Reader
	client.Writer
	client.StatusClient
}

func (c *fakeClient) Get(ctx context.Context, key types.NamespacedName, obj runtime.Object) error {
	return nil
}

func TestSwitchPortStateMachine(t *testing.T) {
	r := SwitchPortReconciler{}
	instance := v1alpha1.SwitchPort{}
	instance.OwnerReferences = []metav1.OwnerReference{
		{
			Kind: "Test",
		},
	}

	m := machine.New(
		&machine.Information{
			Client: &fakeClient{},
			Logger: r.Log,
		},
		&instance,
		&machine.Handlers{
			v1alpha1.SwitchPortNone:        r.noneHandler,
			v1alpha1.SwitchPortIdle:        r.idleHandler,
			v1alpha1.SwitchPortValidating:  r.validatingandler,
			v1alpha1.SwitchPortConfiguring: r.configuringHandler,
			v1alpha1.SwitchPortActive:      r.activeHandler,
			v1alpha1.SwitchPortCleaning:    r.cleaningHandler,
			v1alpha1.SwitchPortDeleting:    r.deletingHandler,
		},
	)

	cases := []struct {
		name                   string
		configurationRef       *v1alpha1.SwitchPortConfigurationRef
		deletionTimestampExist bool
		expectedState          machine.StateType
	}{
		// Delete when `Idle` state
		{
			name:             "<None> -> Idle",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortIdle,
		},
		{
			name:             "Idle -> Validating",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortValidating,
		},
		{
			name:             "Validating -> Configuring",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortConfiguring,
		},
		{
			name:             "Configuring -> Active",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortActive,
		},
		{
			name:          "Active -> Cleaning",
			expectedState: v1alpha1.SwitchPortCleaning,
		},
		{
			name:             "Cleaning -> Idle",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortIdle,
		},
		{
			name:                   "Idle -> Deleting",
			deletionTimestampExist: true,
			expectedState:          v1alpha1.SwitchPortDeleting,
		},
		{
			name:                   "Deleting -> <None>",
			deletionTimestampExist: true,
			expectedState:          v1alpha1.SwitchPortNone,
		},
		// Delete when Cleaning state
		{
			name:             "<None> -> Idle",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortIdle,
		},
		{
			name:             "Idle -> Validating",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortValidating,
		},
		{
			name:             "Validating -> Configuring",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortConfiguring,
		},
		{
			name:             "Configuring -> Active",
			configurationRef: &v1alpha1.SwitchPortConfigurationRef{},
			expectedState:    v1alpha1.SwitchPortActive,
		},
		{
			name:          "Active -> Cleaning",
			expectedState: v1alpha1.SwitchPortCleaning,
		},
		{
			name:                   "Cleaning -> Deleting",
			deletionTimestampExist: true,
			expectedState:          v1alpha1.SwitchPortDeleting,
		},
		{
			name:                   "Deleting -> <None>",
			deletionTimestampExist: true,
			expectedState:          v1alpha1.SwitchPortNone,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			instance.Spec.ConfigurationRef = c.configurationRef
			if c.deletionTimestampExist {
				now := metav1.Now()
				instance.DeletionTimestamp = &now
			} else {
				instance.DeletionTimestamp = nil
			}

			m.Reconcile(context.TODO())
			if c.expectedState != instance.GetState() {
				t.Errorf("Expected state: %s, got: %s", c.expectedState, instance.GetState())
			}
		})
	}
}
