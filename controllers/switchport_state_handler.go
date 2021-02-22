package controllers

import (
	"context"
	"time"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device"
	"github.com/metal3-io/networkconfiguration-operator/pkg/machine"
	"github.com/metal3-io/networkconfiguration-operator/pkg/util/finalizer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const finalizerKey string = "metal3.io.v1alpha1"

// noneHandler add finalizers to CR
func (r *SwitchPortReconciler) noneHandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Add finalizer
	err := finalizer.AddFinalizer(&i.Finalizers, finalizerKey)
	result := reconcile.Result{}
	if err == nil {
		result.Requeue = true
	}

	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, err
}

// idleHandler check spec.configurationRef's value, if isn't nil set the state of CR to `Configuring`
func (r *SwitchPortReconciler) idleHandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchPortDeleting, ctrl.Result{Requeue: true}, nil
	}

	if i.Spec.ConfigurationRef == nil {
		return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, nil
	}

	_, err := i.Spec.ConfigurationRef.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	return v1alpha1.SwitchPortValidating, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchPortReconciler) validatingandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {

	// TODO: Check connection to switch

	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// configuringHandler configure port's network and check configuration progress. If finished set the state of CR to `Configured` state
func (r *SwitchPortReconciler) configuringHandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if i.Spec.ConfigurationRef == nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true}, nil
	}

	dev, err := device.New(ctx, info.Client, &i.OwnerReferences[0])
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	configuration, err := i.Spec.ConfigurationRef.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	err = dev.ConfigurePort(ctx, configuration, i.Spec.ID)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true}, nil
}

// activeHandler check whether the target configuration is consistent with the actual configuration,
// return to `Configuring` state when inconsistent
func (r *SwitchPortReconciler) activeHandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.ConfigurationRef == nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	dev, err := device.New(ctx, info.Client, &i.OwnerReferences[0])
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	configuration, err := i.Spec.ConfigurationRef.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	isIdentical, err := dev.CheckPortConfigutation(ctx, configuration, i.Spec.ID)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	if isIdentical {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, nil
	}

	i.Status.Configuration = nil
	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// cleaningHandler will be called when deleting network configuration, when finished clean spec.configurationRef and status.configurationRef then set CR's state to `Deleted` state.
func (r *SwitchPortReconciler) cleaningHandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	dev, err := device.New(ctx, info.Client, &i.OwnerReferences[0])
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	err = dev.DeConfigurePort(ctx, i.Spec.ID)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	// User delete CR
	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchPortDeleting, ctrl.Result{Requeue: true}, nil
	}

	i.Status.Configuration = nil
	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, err
}

// deletingHandler will remove finalizers, if spec.configurationRef isn't nil set CR's state to <none>
func (r *SwitchPortReconciler) deletingHandler(ctx context.Context, info *machine.Information, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Remove finalizer
	err := finalizer.RemoveFinalizer(&i.Finalizers, finalizerKey)
	result := reconcile.Result{}
	if err != nil {
		result.Requeue = true
	}

	return v1alpha1.SwitchPortNone, result, err
}
