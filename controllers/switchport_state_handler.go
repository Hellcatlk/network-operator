package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/backends"
	"github.com/Hellcatlk/network-operator/pkg/backends/switches"
	"github.com/Hellcatlk/network-operator/pkg/machine"
	"github.com/Hellcatlk/network-operator/pkg/provider"
	"github.com/Hellcatlk/network-operator/pkg/utils/finalizer"
	"k8s.io/apimachinery/pkg/api/errors"
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

	return machine.ResultContinue(v1alpha1.SwitchPortIdle, 0, nil)
}

// idleHandler check spec.configurationRef's value, if isn't nil set the state of CR to `Validating`
func (r *SwitchPortReconciler) idleHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchPortDeleting, 0, nil)
	}

	if i.Spec.Configuration == nil || len(i.OwnerReferences) == 0 {
		return machine.ResultComplete(v1alpha1.SwitchPortIdle, nil)
	}

	return machine.ResultContinue(v1alpha1.SwitchPortVerifying, 0, nil)
}

// verifyingHandler verify the configuration meets the requirements of the switch or not
func (r *SwitchPortReconciler) verifyingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return machine.ResultContinue(v1alpha1.SwitchPortIdle, 0, nil)
	}

	// Fetch switch
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}

	// Fetch configuration
	configuration, err := i.Spec.Configuration.Fetch(ctx, info.Client)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}

	// Check switch limit
	err = owner.Spec.Limit.VerifyConfiguration(configuration)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}

	// Check switch port limit
	err = owner.Status.Ports[i.Name].Verify(configuration)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}

	// Check user limit
	resourceLimit, err := i.FetchSwitchResourceLimit(ctx, info.Client)
	if err != nil && !errors.IsNotFound(err) {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}
	if resourceLimit != nil {
		resource, err := resourceLimit.FetchSwitchResource(ctx, info.Client)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
		}
		for _, limit := range resource.Status.TenantLimits {
			if limit.Namespace == resourceLimit.Namespace {
				err = limit.VerifyConfiguration(configuration)
				if err != nil {
					return machine.ResultContinue(v1alpha1.SwitchPortVerifying,
						requeueAfterTime,
						fmt.Errorf("%s, %s", err, "please check `SwitchResourceLimit/user-limit`"),
					)
				}
			}
		}

	}

	// Check connection with switch
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}
	err = backend.IsAvailable()
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortVerifying, requeueAfterTime, err)
	}

	// Copy configuration to Status.Configuration
	i.Status.Configuration = &configuration.Spec
	i.Status.PortName = owner.Status.Ports[i.Name].Name
	return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, 0, nil)
}

// configuringHandler configure port's network and check configuration progress. If finished set the state of CR to `Active` state
func (r *SwitchPortReconciler) configuringHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return machine.ResultContinue(v1alpha1.SwitchPortCleaning, 0, nil)
	}

	resourceLimit, err := i.FetchSwitchResourceLimit(ctx, info.Client)
	if err != nil && !errors.IsNotFound(err) {
		return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
	}

	// Set configuration to port
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
	}
	err = backend.SetPortAttr(ctx, i.Status.PortName, i.Status.Configuration)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
	}

	if resourceLimit.GetName() != "" {
		err = resourceLimit.Expansion(i.Status.Configuration)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
		}
		err = info.Client.Status().Update(ctx, resourceLimit)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
		}
	}

	return machine.ResultContinue(v1alpha1.SwitchPortActive, 0, nil)
}

// activeHandler check whether the target configuration is consistent with the actual configuration,
// return to `Configuring` state when inconsistent
func (r *SwitchPortReconciler) activeHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	if !i.DeletionTimestamp.IsZero() || i.Spec.Configuration == nil {
		return machine.ResultContinue(v1alpha1.SwitchPortCleaning, 0, nil)
	}

	// Check spec.ConfigurationRef as same as status.Configuration or not
	configuration, err := i.Spec.Configuration.Fetch(ctx, info.Client)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortActive, requeueAfterTime, err)
	}
	if !i.Status.Configuration.IsEqual(&configuration.Spec) {
		return machine.ResultContinue(v1alpha1.SwitchPortCleaning, 0, nil)
	}

	// Check status.Configuration as same as switch's port configuration or not
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortActive, requeueAfterTime, err)
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortActive, requeueAfterTime, err)
	}
	actualConfiguration, err := backend.GetPortAttr(ctx, i.Status.PortName)
	if err != nil || i.Status.Configuration.IsEqual(actualConfiguration) {
		return machine.ResultContinue(v1alpha1.SwitchPortActive, requeueAfterTime, err)
	}

	info.Logger.Info("configuration of port has been changed externally")
	return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, 0, nil)
}

// cleaningHandler will be called when deleting network configuration, when finished clean spec.configurationRef and status.configurationRef then set CR's state to `Idle` state.
func (r *SwitchPortReconciler) cleaningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)
	resourceLimit, err := i.FetchSwitchResourceLimit(ctx, info.Client)
	if err != nil && !errors.IsNotFound(err) {
		return machine.ResultContinue(v1alpha1.SwitchPortConfiguring, requeueAfterTime, err)
	}
	// Remove switch's port configuration
	owner, err := i.FetchOwnerReference(ctx, info.Client)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortCleaning, requeueAfterTime, err)
	}
	backend, err := getSwitchBackend(ctx, info.Client, owner)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortCleaning, requeueAfterTime, err)
	}
	err = backend.ResetPort(ctx, i.Status.PortName, i.Status.Configuration)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchPortCleaning, requeueAfterTime, err)
	}

	if resourceLimit.GetName() != "" {
		err = resourceLimit.Shrink(i.Status.Configuration)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchPortCleaning, requeueAfterTime, err)
		}
		err = info.Client.Status().Update(ctx, resourceLimit)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchPortCleaning, requeueAfterTime, err)
		}
	}
	i.Status.Configuration = nil
	i.Status.PortName = ""
	return machine.ResultContinue(v1alpha1.SwitchPortIdle, 0, nil)
}

// deletingHandler will remove finalizers
func (r *SwitchPortReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchPort)

	// Remove finalizer
	finalizer.Remove(&i.Finalizers, finalizerKey)

	return machine.ResultComplete(v1alpha1.SwitchPortDeleting, nil)
}
