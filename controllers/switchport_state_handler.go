package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Hellcatlk/networkconfiguration-operator/api/v1alpha1"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/devices/switches"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/machine"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/provider"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/utils/finalizer"
	ctrl "sigs.k8s.io/controller-runtime"
)

const switchPortFinalizerKey string = "metal3.io.switchport"
const requeueAfterTime time.Duration = time.Second * 10

// noneHandler add finalizers to CR
func (r *SwitchPortReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("none")

	i := instance.(*v1alpha1.SwitchPort)

	// Add finalizer
	finalizer.Add(&i.Finalizers, switchPortFinalizerKey)

	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, nil
}

// idleHandler check spec.configurationRef's value, if isn't nil set the state of CR to `Validating`
func (r *SwitchPortReconciler) idleHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("idle")

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
	info.Logger.Info("verifying")

	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, nil
	}

	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	configuration, err := i.Spec.Configuration.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	if owner.Spec.Ports[i.Name].Disabled {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, fmt.Errorf("the port is disabled")
	}

	if owner.Spec.Ports[i.Name].TrunkDisabled && len(configuration.Spec.Vlans) != 0 {
		return v1alpha1.SwitchPortVerifying, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, fmt.Errorf("set the port to trunk mode is disabled")
	}

	// Copy configuration to Status.Configuration
	i.Status.Configuration = configuration
	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// configuringHandler configure port's network and check configuration progress. If finished set the state of CR to `Active` state
func (r *SwitchPortReconciler) configuringHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("configing")

	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true}, nil
	}

	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	providerSwitch, err := owner.Spec.ProviderSwitch.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	config, err := provider.GetSwitchConfiguration(ctx, info.Client, providerSwitch)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	sw, err := switches.New(ctx, config)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Set configuration to port
	err = sw.SetPortAttr(ctx, owner.Spec.Ports[i.Name].Name, i.Status.Configuration)
	if err != nil {
		return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true}, nil
}

// activeHandler check whether the target configuration is consistent with the actual configuration,
// return to `Configuring` state when inconsistent
func (r *SwitchPortReconciler) activeHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("active")

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

	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	providerSwitch, err := owner.Spec.ProviderSwitch.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	config, err := provider.GetSwitchConfiguration(ctx, info.Client, providerSwitch)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	sw, err := switches.New(ctx, config)
	if err != nil {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Check status.Configuration as same as switch's port configuration or not
	configuration, err = sw.GetPortAttr(ctx, owner.Spec.Ports[i.Name].Name)
	if err != nil || reflect.DeepEqual(configuration.Spec, i.Status.Configuration.Spec) {
		return v1alpha1.SwitchPortActive, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	info.Logger.Info("configuration of port has been changed externally")
	return v1alpha1.SwitchPortConfiguring, ctrl.Result{Requeue: true}, nil
}

// cleaningHandler will be called when deleting network configuration, when finished clean spec.configurationRef and status.configurationRef then set CR's state to `Idle` state.
func (r *SwitchPortReconciler) cleaningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("cleaning")

	i := instance.(*v1alpha1.SwitchPort)

	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	providerSwitch, err := owner.Spec.ProviderSwitch.Fetch(ctx, info.Client)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	config, err := provider.GetSwitchConfiguration(ctx, info.Client, providerSwitch)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	sw, err := switches.New(ctx, config)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Remove switch's port configuration
	err = sw.ResetPort(ctx, owner.Spec.Ports[i.Name].Name, i.Status.Configuration)
	if err != nil {
		return v1alpha1.SwitchPortCleaning, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	i.Status.Configuration = nil
	return v1alpha1.SwitchPortIdle, ctrl.Result{Requeue: true}, err
}

// deletingHandler will remove finalizers
func (r *SwitchPortReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("deleting")

	i := instance.(*v1alpha1.SwitchPort)

	// Remove finalizer
	finalizer.Remove(&i.Finalizers, switchPortFinalizerKey)

	return v1alpha1.SwitchPortNone, ctrl.Result{}, nil
}
