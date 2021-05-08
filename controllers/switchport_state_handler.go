package controllers

import (
	"context"
	"reflect"
	"time"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/device/switches"
	"github.com/metal3-io/networkconfiguration-operator/pkg/machine"
	"github.com/metal3-io/networkconfiguration-operator/pkg/utils/finalizer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const finalizerKey string = "metal3.io.v1alpha1"
const requeueAfterTime time.Duration = time.Second * 10

// noneHandler add finalizers to CR
func (r *SwitchPortReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Add finalizer
	err := finalizer.Add(&i.Finalizers, finalizerKey)
	result := reconcile.Result{}
	if err == nil {
		result.Requeue = true
	}

	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, err
}

// idleHandler check spec.configurationRef's value, if isn't nil set the state of CR to `Validating`
func (r *SwitchPortReconciler) idleHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchPortDeleting, ctrl.Result{Requeue: true}, nil
	}

	if i.Spec.ConfigurationRef == nil || len(i.OwnerReferences) == 0 {
		return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, nil
	}

	return v1alpha1.SwitchPortValidating, ctrl.Result{Requeue: true}, nil
}

func (r *SwitchPortReconciler) validatingandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// TODO: Check connection with switch

	// Copy configuration to Status.Configuration
	configuration, err := i.Spec.ConfigurationRef.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	i.Status.Configuration = configuration

	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// configuringHandler configure port's network and check configuration progress. If finished set the state of CR to `Active` state
func (r *SwitchPortReconciler) configuringHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.ConfigurationRef == nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	// Set port configuration to switch
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	sw, err := switches.New(ctx, owner.Spec.OS, owner.Spec.URL, owner.Spec.Username, owner.Spec.Password, owner.Spec.Options)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	err = sw.SetPortAttr(ctx, i.Spec.ID, i.Status.Configuration)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true}, nil
}

// activeHandler check whether the target configuration is consistent with the actual configuration,
// return to `Configuring` state when inconsistent
func (r *SwitchPortReconciler) activeHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.ConfigurationRef == nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	// Check spec.ConfigurationRef as same as status.Configuration or not
	configuration, err := i.Spec.ConfigurationRef.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	if !reflect.DeepEqual(configuration, i.Status.Configuration) {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	// Check status.Configuration as same as switch's port configuration or not
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	sw, err := switches.New(ctx, owner.Spec.OS, owner.Spec.URL, owner.Spec.Username, owner.Spec.Password, owner.Spec.Options)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	configuration, err = sw.GetPortAttr(ctx, i.Spec.ID)
	if err != nil || reflect.DeepEqual(configuration, i.Status.Configuration) {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// cleaningHandler will be called when deleting network configuration, when finished clean spec.configurationRef and status.configurationRef then set CR's state to `Idle` state.
func (r *SwitchPortReconciler) cleaningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Remove switch's port configuration
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	sw, err := switches.New(ctx, owner.Spec.OS, owner.Spec.URL, owner.Spec.Username, owner.Spec.Password, owner.Spec.Options)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	err = sw.ResetPort(ctx, i.Spec.ID, i.Status.Configuration)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	i.Status.Configuration = nil
	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, err
}

// deletingHandler will remove finalizers
func (r *SwitchPortReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	result := reconcile.Result{}
	// Remove finalizer
	err := finalizer.Remove(&i.Finalizers, finalizerKey)
	if err != nil {
		result.Requeue = true
	}

	return v1alpha1.SwitchPortNone, result, err
}
