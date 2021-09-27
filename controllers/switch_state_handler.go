package controllers

import (
	"context"
	"fmt"
	"reflect"

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
	i := instance.(*v1alpha1.Switch)

	// Add finalizer
	finalizer.Add(&i.Finalizers, finalizerKey)

	return machine.ResultContinue(v1alpha1.SwitchVerifying, 0, nil)
}

func (r *SwitchReconciler) verifyingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchDeleting, 0, nil)
	}

	// Check connection with switch
	backend, err := getSwitchBackend(ctx, info.Client, i)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchVerifying, requeueAfterTime, err)
	}
	err = backend.IsAvailable()
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchVerifying, requeueAfterTime, err)
	}

	// Delete SwitchPorts which isn't included i.Spec
	for name := range i.Status.Ports {
		_, exist := i.Spec.Ports[name]
		if !exist || !reflect.DeepEqual(i.Spec.Ports[name], i.Status.Ports[name]) {
			switchPort := &v1alpha1.SwitchPort{}
			switchPort.Name = name
			switchPort.Namespace = i.Namespace
			err := info.Client.Delete(ctx, switchPort)
			if err != nil {
				if errors.IsNotFound(err) {
					continue
				}
				return machine.ResultContinue(v1alpha1.SwitchVerifying, requeueAfterTime, err)
			}
		}
	}

	if i.Status.Provider == nil {
		if i.Spec.Provider == nil || reflect.DeepEqual(i.Spec.Provider, &v1alpha1.SwitchProviderReference{}) {
			return machine.ResultContinue(v1alpha1.SwitchVerifying, requeueAfterTime, fmt.Errorf("provider is nil or empty"))
		}
		i.Status.Provider = i.Spec.Provider.DeepCopy()
	} else {
		info.Logger.Info("the provider field is not allowed to be edited")
	}
	i.Status.Ports = i.Spec.Ports
	return machine.ResultContinue(v1alpha1.SwitchConfiguring, 0, nil)
}

func (r *SwitchReconciler) configuringHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchDeleting, 0, nil)
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
		if err != nil {
			// If SwitchPort is existed, skip it
			if errors.IsAlreadyExists(err) {
				continue
			}
			return machine.ResultContinue(v1alpha1.SwitchConfiguring, requeueAfterTime, err)
		}
	}

	return machine.ResultContinue(v1alpha1.SwitchRunning, 0, nil)
}

func (r *SwitchReconciler) runningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Switch)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchDeleting, 0, nil)
	}

	if !reflect.DeepEqual(i.Spec.Ports, i.Status.Ports) {
		return machine.ResultContinue(v1alpha1.SwitchVerifying, 0, nil)
	}

	// Check SwitchPorts are existed
	for name := range i.Status.Ports {
		// Get SwitchPort
		err := info.Client.Get(
			ctx, types.NamespacedName{
				Name:      name,
				Namespace: i.Namespace,
			},
			&v1alpha1.SwitchPort{},
		)
		if err != nil {
			// If SwitchPort isn't find, return configuring state and create it
			if errors.IsNotFound(err) {
				return machine.ResultContinue(v1alpha1.SwitchConfiguring, 0, err)
			}
			return machine.ResultContinue(v1alpha1.SwitchRunning, requeueAfterTime, err)
		}
	}

	// Check connection with switch
	backend, err := getSwitchBackend(ctx, info.Client, i)
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchRunning, requeueAfterTime, err)
	}
	err = backend.IsAvailable()
	if err != nil {
		return machine.ResultContinue(v1alpha1.SwitchRunning, requeueAfterTime, err)
	}

	return machine.ResultContinue(v1alpha1.SwitchRunning, requeueAfterTime, nil)
}

func (r *SwitchReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.Switch)

	// Foreground delete
	propagationPolicy := metav1.DeletePropagationForeground
	err := info.Client.Delete(ctx, i, &client.DeleteOptions{PropagationPolicy: &propagationPolicy})
	if err != nil {
		return v1alpha1.SwitchDeleting, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Remove finalizer
	finalizer.Remove(&i.Finalizers, finalizerKey)

	return machine.ResultComplete(v1alpha1.SwitchDeleting, nil)
}
