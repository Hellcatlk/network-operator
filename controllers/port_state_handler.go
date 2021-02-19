package controllers

import (
	"context"
	"fmt"
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
func (r *PortReconciler) noneHandler(ctx context.Context, info *machine.Information, instance interface{}) (v1alpha1.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Port)

	// Add finalizer
	err := finalizer.AddFinalizer(&i.Finalizers, finalizerKey)
	result := reconcile.Result{}
	if err == nil {
		result.Requeue = true
	}

	return v1alpha1.PortCreated, ctrl.Result{Requeue: true}, err
}

// createdHandler check spec.configurationRef's value, if isn't nil set the state of CR to `Configuring`
func (r *PortReconciler) createdHandler(ctx context.Context, info *machine.Information, instance interface{}) (v1alpha1.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Port)

	// Initialize device
	dev, err := device.New(ctx, info.Client, &i.OwnerReferences[0])
	if err != nil {
		return v1alpha1.PortCreated, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	// Get port's state
	switch dev.PortState(ctx, i.Spec.ID) {
	case device.None, device.Deleted:
		// Go to `Configuring` state
		return v1alpha1.PortConfiguring, ctrl.Result{Requeue: true}, nil

	case device.Deleting:
		// Just wait

	default:
		return v1alpha1.PortCleaned, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, fmt.Errorf("port(%s) have been used", i.Spec.ID)
	}

	return v1alpha1.PortCreated, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, nil
}

// configuringHandler configure port's network and check configuration progress. If finished set the state of CR to `Configured` state
func (r *PortReconciler) configuringHandler(ctx context.Context, info *machine.Information, instance interface{}) (v1alpha1.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Port)

	dev, err := device.New(ctx, info.Client, &i.OwnerReferences[0])
	if err != nil {
		return v1alpha1.PortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	switch dev.PortState(ctx, i.Spec.ID) {
	case device.None, device.Deleted, device.ConfigureFailed:
		// Fetch network configuration
		configuration, err := i.Spec.ConfigurationRef.Fetch(ctx, info.Client)
		if err != nil {
			return v1alpha1.PortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
		}
		// Configure network
		err = dev.ConfigurePort(ctx, configuration, i.Spec.ID)

	case device.Configured:
		// If configure network success, we just need to set next state to `Configured`, but not reconcile
		return v1alpha1.PortConfigured, ctrl.Result{Requeue: false}, nil

	default:
		// Just wait
	}

	return v1alpha1.PortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
}

// configuredHandler check whether the target configuration is consistent with the actual configuration,
// return to `Configuring` state when inconsistent
func (r *PortReconciler) configuredHandler(ctx context.Context, info *machine.Information, instance interface{}) (v1alpha1.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Port)

	// User update CR
	if i.DeletionTimestamp.IsZero() {
		return v1alpha1.PortConfiguring, ctrl.Result{Requeue: true}, nil
	}

	// User delete CR
	return v1alpha1.PortCleaning, ctrl.Result{Requeue: true}, nil
}

// cleaningHandler will be called when deleting network configuration, when finished clean spec.configurationRef and status.configurationRef then set CR's state to `Deleted` state.
func (r *PortReconciler) cleaningHandler(ctx context.Context, info *machine.Information, instance interface{}) (v1alpha1.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Port)

	dev, err := device.New(ctx, info.Client, &i.OwnerReferences[0])
	if err != nil {
		return v1alpha1.PortCleaning, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
	}

	switch dev.PortState(ctx, i.Spec.ID) {
	case device.Configured, device.ConfigureFailed, device.DeleteFailed:
		// Delete network
		err = dev.DeConfigurePort(ctx, i.Spec.ID)

	case device.None, device.Deleted:
		return v1alpha1.PortCleaned, ctrl.Result{Requeue: true}, nil

	default:
		// Just wait
	}

	return v1alpha1.PortCleaning, ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
}

// cleanedHandler will remove finalizers, if spec.configurationRef isn't nil set CR's state to <none>
func (r *PortReconciler) cleanedHandler(ctx context.Context, info *machine.Information, instance interface{}) (v1alpha1.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Port)

	// Remove finalizer
	err := finalizer.RemoveFinalizer(&i.Finalizers, finalizerKey)
	result := reconcile.Result{}
	if err != nil {
		result.Requeue = true
	}

	return v1alpha1.PortCleaned, result, err
}
