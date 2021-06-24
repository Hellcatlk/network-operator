package controllers

import (
	"context"
	"testing"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type fakeClient struct {
	client.Reader
	client.Writer
	client.StatusClient
}

func (c *fakeClient) Get(ctx context.Context, key types.NamespacedName, obj runtime.Object) error {
	switch key.Name {
	case "Switch":
		*obj.(*v1alpha1.Switch) = v1alpha1.Switch{
			Status: v1alpha1.SwitchStatus{
				ProviderSwitch: &v1alpha1.ProviderSwitchRef{
					Kind: "TestSwitch",
					Name: "TestSwitch",
				},
			},
		}
	case "SwitchPortConfiguration":
		*obj.(*v1alpha1.SwitchPortConfiguration) = v1alpha1.SwitchPortConfiguration{}
	case "Secret":
		*obj.(*corev1.Secret) = corev1.Secret{}
	}

	return nil
}

func TestSwitchPortStateMachine(t *testing.T) {
	r := SwitchPortReconciler{}
	instance := v1alpha1.SwitchPort{}
	instance.OwnerReferences = []metav1.OwnerReference{
		{
			Name: "Switch",
		},
	}

	m := machine.New(
		&machine.ReconcileInfo{
			Client: &fakeClient{},
			Logger: log.NullLogger{},
		},
		&instance,
		&machine.Handlers{
			v1alpha1.SwitchPortNone:        r.noneHandler,
			v1alpha1.SwitchPortIdle:        r.idleHandler,
			v1alpha1.SwitchPortVerifying:   r.verifyingHandler,
			v1alpha1.SwitchPortConfiguring: r.configuringHandler,
			v1alpha1.SwitchPortActive:      r.activeHandler,
			v1alpha1.SwitchPortCleaning:    r.cleaningHandler,
			v1alpha1.SwitchPortDeleting:    r.deletingHandler,
		},
	)

	cases := []struct {
		name                   string
		configurationRefExist  bool
		deletionTimestampExist bool
		expectedDirty          bool
		expectedState          machine.StateType
	}{
		// Delete when `Idle` state
		{
			name:          "<None> -> Idle",
			expectedDirty: true,
			expectedState: v1alpha1.SwitchPortIdle,
		},
		{
			name:          "Idle -> Idle",
			expectedDirty: false,
			expectedState: v1alpha1.SwitchPortIdle,
		},
		{
			name:                  "Idle -> Verifying",
			configurationRefExist: true,
			expectedDirty:         true,
			expectedState:         v1alpha1.SwitchPortVerifying,
		},
		{
			name:                  "Verifying -> Configuring",
			configurationRefExist: true,
			expectedDirty:         true,
			expectedState:         v1alpha1.SwitchPortConfiguring,
		},
		{
			name:                  "Configuring -> Active",
			configurationRefExist: true,
			expectedDirty:         true,
			expectedState:         v1alpha1.SwitchPortActive,
		},
		{
			name:                  "Active -> Active",
			configurationRefExist: true,
			expectedDirty:         false,
			expectedState:         v1alpha1.SwitchPortActive,
		},
		{
			name:          "Active -> Cleaning",
			expectedDirty: true,
			expectedState: v1alpha1.SwitchPortCleaning,
		},
		{
			name:                  "Cleaning -> Idle",
			configurationRefExist: true,
			expectedDirty:         true,
			expectedState:         v1alpha1.SwitchPortIdle,
		},
		{
			name:                   "Idle -> Deleting",
			deletionTimestampExist: true,
			expectedDirty:          true,
			expectedState:          v1alpha1.SwitchPortDeleting,
		},
		{
			name:                   "Deleting -> Deleting",
			deletionTimestampExist: true,
			expectedDirty:          true,
			expectedState:          v1alpha1.SwitchPortDeleting,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.configurationRefExist {
				instance.Spec.Configuration = &v1alpha1.SwitchPortConfigurationReference{
					Name: "SwitchPortConfiguration",
				}
			} else {
				instance.Spec.Configuration = nil
			}

			if c.deletionTimestampExist {
				now := metav1.Now()
				instance.DeletionTimestamp = &now
			} else {
				instance.DeletionTimestamp = nil
			}

			dirty, _, err := m.Reconcile(context.TODO())
			if c.expectedDirty != dirty {
				t.Errorf("Expected dirty: %v, got: %v", c.expectedDirty, dirty)
			}
			if c.expectedState != instance.GetState() {
				t.Errorf("Expected state: %s, got: %s", c.expectedState, instance.GetState())
			}
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
