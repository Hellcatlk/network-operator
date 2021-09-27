package controllers

import (
	"context"
	"testing"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	"github.com/Hellcatlk/network-operator/pkg/utils/finalizer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestSwitchStateMachine(t *testing.T) {
	r := SwitchReconciler{}
	instance := v1alpha1.Switch{}
	instance.Name = "Switch"
	instance.Spec.Provider = &v1alpha1.SwitchProviderReference{
		Kind: "TestSwitch",
		Name: "TestSwitch",
	}

	m := machine.New(
		&machine.ReconcileInfo{
			Client: &fakeClient{},
			Logger: log.NullLogger{},
		},
		&instance,
		map[machine.StateType]machine.Handler{
			v1alpha1.SwitchNone:        r.noneHandler,
			v1alpha1.SwitchVerifying:   r.verifyingHandler,
			v1alpha1.SwitchConfiguring: r.configuringHandler,
			v1alpha1.SwitchRunning:     r.runningHandler,
			v1alpha1.SwitchDeleting:    r.deletingHandler,
		},
	)

	cases := []struct {
		name                   string
		finalizerExist         bool
		deletionTimestampExist bool
		expectedDirty          bool
		expectedState          machine.StateType
		expectedError          bool
	}{
		{
			name:           "<None> -> Verifying",
			finalizerExist: true,
			expectedDirty:  true,
			expectedState:  v1alpha1.SwitchVerifying,
		},
		{
			name:          "Verifying -> Configuring",
			expectedDirty: true,
			expectedState: v1alpha1.SwitchConfiguring,
		},
		{
			name:          "Configuring -> Running",
			expectedDirty: true,
			expectedState: v1alpha1.SwitchRunning,
		},
		{
			name:                   "Running -> Deleting",
			expectedDirty:          true,
			deletionTimestampExist: true,
			expectedState:          v1alpha1.SwitchDeleting,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.finalizerExist {
				finalizer.Add(&instance.Finalizers, finalizerKey)
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
			if c.expectedError != (err != nil) {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
