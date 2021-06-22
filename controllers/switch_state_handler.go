package controllers

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	"github.com/Hellcatlk/network-operator/pkg/utils/finalizer"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// noneHandler add finalizers to CR
func (r *SwitchReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("none")

	i := instance.(*v1alpha1.Switch)

	// Add finalizer
	finalizer.Add(&i.Finalizers, finalizerKey)

	return v1alpha1.SwitchVerify, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchReconciler) verifyingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("verifying")

	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchDeleting, ctrl.Result{Requeue: true}, nil
	}

	i.Status.ProviderSwitch = i.Spec.ProviderSwitch.DeepCopy()
	i.Status.Ports = i.Spec.Ports
	return v1alpha1.SwitchConfiguring, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchReconciler) configuringHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("configuring")

	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchDeleting, ctrl.Result{Requeue: true}, nil
	}

	// Create SwitchPorts
	for name := range i.Status.Ports {
		switchPort := &v1alpha1.SwitchPort{}
		switchPort.Name = name
		switchPort.Namespace = i.Namespace
		switchPort.OwnerReferences = []metav1.OwnerReference{
			{
				BlockOwnerDeletion: new(bool),
				APIVersion:         i.APIVersion,
				Kind:               i.Kind,
				Name:               i.Name,
				UID:                i.UID,
			},
		}
		*switchPort.OwnerReferences[0].BlockOwnerDeletion = true

		// Create SwitchPort
		err := info.Client.Create(ctx, switchPort)
		if !errors.IsAlreadyExists(err) {
			return v1alpha1.SwitchConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
		}
	}

	return v1alpha1.SwitchRunning, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchReconciler) runningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("running")

	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchDeleting, ctrl.Result{Requeue: true}, nil
	}

	// Check SwitchPorts are existed
	for name := range i.Status.Ports {
		err := info.Client.Get(
			ctx, types.NamespacedName{
				Name:      name,
				Namespace: i.Namespace,
			},
			&v1alpha1.SwitchPort{},
		)
		if errors.IsNotFound(err) {
			return v1alpha1.SwitchConfiguring, ctrl.Result{Requeue: true}, nil
		}
	}

	return v1alpha1.SwitchRunning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, nil
}

func (r *SwitchReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("deleting")

	i := instance.(*v1alpha1.Switch)

	// Remove finalizer
	finalizer.Remove(&i.Finalizers, finalizerKey)

	// Set foreground delete policy
	propagationPolicy := metav1.DeletePropagationForeground
	info.Client.Delete(ctx, i, &client.DeleteOptions{PropagationPolicy: &propagationPolicy})

	return v1alpha1.SwitchDeleting, ctrl.Result{}, nil
}
