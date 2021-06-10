package controllers

import (
	"context"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/devices/switches"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	"github.com/Hellcatlk/network-operator/pkg/utils/finalizer"
	ctrl "sigs.k8s.io/controller-runtime"
)

const switchFinalizerKey string = "foregroundDeletion"

// noneHandler add finalizers to CR
func (r *SwitchReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("none")

	i := instance.(*v1alpha1.Switch)

	// Add finalizer
	finalizer.Add(&i.Finalizers, switchFinalizerKey)

	return v1alpha1.SwitchVerify, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchReconciler) verifyingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("verifying")

	i := instance.(*v1alpha1.Switch)

	providerSwitch, err := i.Spec.ProviderSwitch.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchVerify, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	config, err := providerSwitch.GetConfiguration(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchVerify, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	sw, err := switches.New(ctx, config)
	if err != nil {
		return v1alpha1.SwitchVerify, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	err = sw.PowerOn(ctx)
	if err != nil {
		return v1alpha1.SwitchVerify, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	i.Status.ProviderSwitch = i.Spec.ProviderSwitch.DeepCopy()
	i.Status.Ports = i.Spec.Ports
	return v1alpha1.SwitchCreating, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchReconciler) creatingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {

	return v1alpha1.SwitchActive, ctrl.Result{}, nil
}

func (r *SwitchReconciler) activeHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchDeleting, ctrl.Result{}, nil
	}

	return v1alpha1.SwitchActive, ctrl.Result{}, nil
}

func (r *SwitchReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {

	return v1alpha1.SwitchDeleting, ctrl.Result{}, nil
}
