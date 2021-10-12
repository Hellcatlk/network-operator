package controllers

import (
	"context"
	"testing"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestSwitchResourceStateMachine(t *testing.T) {
	r := SwitchResourceReconciler{}
	instance := v1alpha1.SwitchResource{
		Spec: v1alpha1.SwitchResourceSpec{
			VLANRange: "1-1000",
			TenantLimits: map[string]*v1alpha1.TenantLimit{
				"user1": {
					Namespace: "test1",
					VLANRange: "1-10",
				},
			},
		},
	}
	instance.Name = "SwitchResource"

	m := machine.New(
		&machine.ReconcileInfo{
			Client: &fakeClient{},
			Logger: log.NullLogger{},
		},
		&instance,
		map[machine.StateType]machine.Handler{
			v1alpha1.SwitchResourceNone:      r.noneHandler,
			v1alpha1.SwitchResourceVerifying: r.verifyingHandler,
			v1alpha1.SwitchResourceCreating:  r.creatingHandler,
			v1alpha1.SwitchResourceRunning:   r.runningHandler,
			v1alpha1.SwitchResourceDeleting:  r.deletingHandler,
		},
	)

	cases := []struct {
		name                   string
		deletionTimestampExist bool
		expectedDirty          machine.DirtyType
		expectedState          machine.StateType
		expectedError          bool
	}{
		{
			name:          "<None> -> Verifying",
			expectedDirty: machine.All,
			expectedState: v1alpha1.SwitchResourceVerifying,
		},
		{
			name:          "Verifying -> Configuring",
			expectedDirty: machine.Status,
			expectedState: v1alpha1.SwitchResourceCreating,
		},
		{
			name:          "Configuring -> Running",
			expectedDirty: machine.Status,
			expectedState: v1alpha1.SwitchResourceRunning,
		},
		{
			name:          "Running -> Running",
			expectedDirty: machine.None,
			expectedState: v1alpha1.SwitchResourceRunning,
		},
		{
			name:                   "Running -> Deleting",
			deletionTimestampExist: true,
			expectedDirty:          machine.Status,
			expectedState:          v1alpha1.SwitchResourceDeleting,
		},
		{
			name:                   "Deleting -> Deleting",
			deletionTimestampExist: true,
			expectedDirty:          machine.MetadataAndSpec,
			expectedState:          v1alpha1.SwitchResourceDeleting,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
			if c.expectedError != (err != nil) {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
