package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/backends/switches"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/finalizer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const finalizerKey string = "metal3.io"
const requeueAfterTime time.Duration = time.Second * 10

// getSwitchBackend return switch backend
func getSwitchBackend(ctx context.Context, client client.Client, sw *v1alpha1.Switch) (backends.Switch, error) {
	var provider provider.Switch
	var err error
	if sw.Status.Provider != nil {
		provider, err = sw.Status.Provider.Fetch(ctx, client)
		if err != nil {
			return nil, err
		}
	} else {
		provider, err = sw.Spec.Provider.Fetch(ctx, client)
		if err != nil {
			return nil, err
		}
	}
	if provider == nil {
		return nil, fmt.Errorf("can not fetch provider")
	}

	config, err := provider.GetConfiguration(ctx, client)
	if err != nil {
		return nil, err
	}

	return switches.New(ctx, config)
}

// noneHandler add finalizers to CR
func (r *SwitchPortReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Add finalizer
	finalizer.Add(&i.Finalizers, finalizerKey)

	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, nil
}

// idleHandler check spec.configurationRef's value, if isn't nil set the state of CR to `Validating`
func (r *SwitchPortReconciler) idleHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() {
		return v1alpha1.SwitchPortDeleting, ctrl.Result{Requeue: true}, nil
	}

	if i.Spec.Configuration == nil || len(i.OwnerReferences) == 0 {
		return v1alpha1.SwitchPortIdle, ctrl.Result{}, nil
	}

	return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true}, nil
}

// verifyingHandler verify the configuration meets the requirements of the switch or not
func (r *SwitchPortReconciler) verifyingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, nil
	}

	// Check connection with switch
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	err = backend.IsAvaliable()
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Check switch port configuration
	configuration, err := i.Spec.Configuration.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	err = owner.Status.Ports[i.Name].Verify(configuration)
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Copy configuration to Status.Configuration
	i.Status.Configuration = configuration
	i.Status.PortName = owner.Status.Ports[i.Name].Name
	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// configuringHandler configure port's network and check configuration progress. If finished set the state of CR to `Active` state
func (r *SwitchPortReconciler) configuringHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	// Set configuration to port
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	err = backend.SetPortAttr(ctx, i.Status.PortName, i.Status.Configuration)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true}, nil
}

// activeHandler check whether the target configuration is consistent with the actual configuration,
// return to `Configuring` state when inconsistent
func (r *SwitchPortReconciler) activeHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	// Check spec.ConfigurationRef as same as status.Configuration or not
	configuration, err := i.Spec.Configuration.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	if !reflect.DeepEqual(configuration.Spec, i.Status.Configuration.Spec) {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	// Check status.Configuration as same as switch's port configuration or not
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	configuration, err = backend.GetPortAttr(ctx, i.Status.PortName)
	if err != nil || reflect.DeepEqual(configuration.Spec, i.Status.Configuration.Spec) {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	info.Logger.Info("configuration of port has been changed externally")
	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// cleaningHandler will be called when deleting network configuration, when finished clean spec.configurationRef and status.configurationRef then set CR's state to `Idle` state.
func (r *SwitchPortReconciler) cleaningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Remove switch's port configuration
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	err = backend.ResetPort(ctx, i.Status.PortName, i.Status.Configuration)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	i.Status.Configuration = nil
	i.Status.PortName = ""
	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, err
}

// deletingHandler will remove finalizers
func (r *SwitchPortReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Remove finalizer
	finalizer.Remove(&i.Finalizers, finalizerKey)

	return v1alpha1.SwitchPortDeleting, ctrl.Result{}, nil
}
