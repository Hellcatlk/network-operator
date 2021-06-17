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
)

const foregroundDeletionFinalizerKey string = "foregroundDeletion"

// noneHandler add finalizers to CR
func (r *SwitchReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("none")

	i := instance.(*v1alpha1.Switch)

	// Add finalizer
	finalizer.Add(&i.Finalizers, finalizerKey)
	finalizer.Add(&i.Finalizers, foregroundDeletionFinalizerKey)

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
				APIVersion: i.APIVersion,
				Kind:       i.Kind,
				Name:       i.Name,
				UID:        i.UID,
			},
		}

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

	// Check all SwitchPorts have been deleted
	for name := range i.Status.Ports {
		err := info.Client.Get(
			ctx, types.NamespacedName{
				Name:      name,
				Namespace: i.Namespace,
			},
			&v1alpha1.SwitchPort{},
		)
		if err == nil || !errors.IsNotFound(err) {
			return v1alpha1.SwitchDeleting, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, nil
		}
	}

	finalizer.Remove(&i.Finalizers, finalizerKey)

	return v1alpha1.SwitchNone, ctrl.Result{}, nil
}
