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
func (r *SwitchResourceReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchResource)

	// Add finalizer
	finalizer.Add(&i.Finalizers, finalizerKey)

	i.Status.AvailableVLAN = i.Spec.VLANRange

	return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, 0, nil)
}

func (r *SwitchResourceReconciler) verifyingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchResource)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchResourceDeleting, 0, nil)
	}

	// Delete SwitchResourceLimit which isn't included i.Spec
	for name, limit := range i.Status.TenantLimits {
		_, exit := i.Spec.TenantLimits[name]
		if !exit {
			sr, err := limit.FetchSwitchResourceLimit(ctx, info.Client)
			if err != nil {
				return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, requeueAfterTime, err)
			}
			if sr.Status.UsedVLAN != "" {
				err = fmt.Errorf("SwitchResourceLimit %s still has vlan %s being used and cannot be deleted", sr.Name, sr.Status.UsedVLAN)
				return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, requeueAfterTime, err)
			}
		}
		if !exit || !reflect.DeepEqual(i.Spec.TenantLimits[name], i.Status.TenantLimits[name]) {
			switchResourceLimit := &v1alpha1.SwitchResourceLimit{}
			switchResourceLimit.Name = "user-limit"
			switchResourceLimit.Namespace = limit.Namespace
			err := info.Client.Delete(ctx, switchResourceLimit)
			if err != nil {
				if errors.IsNotFound(err) {
					continue
				}
				return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, requeueAfterTime, err)
			}
			err = i.Expansion(limit)
			if err != nil {
				return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, 0, err)
			}
		}
	}

	for _, limit := range i.Spec.TenantLimits {
		err := limit.Verify(&i.Status)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, 0, err)
		}
	}

	i.Status.TenantLimits = i.Spec.TenantLimits

	return machine.ResultContinue(v1alpha1.SwitchResourceCreating, 0, nil)
}

func (r *SwitchResourceReconciler) creatingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchResource)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchDeleting, 0, nil)
	}

	// Create SwitchResourceLimit
	for _, limit := range i.Spec.TenantLimits {
		// Get switchResourceLimit
		err := info.Client.Get(
			ctx, types.NamespacedName{
				Name:      "user-limit",
				Namespace: limit.Namespace,
			},
			&v1alpha1.SwitchResourceLimit{},
		)
		if err == nil {
			continue
		}
		if err != nil {
			if !errors.IsNotFound(err) {
				return machine.ResultContinue(v1alpha1.SwitchResourceCreating, 0, err)
			}
		}

		err = limit.Verify(&i.Status)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchResourceCreating, 0, err)
		}
		switchResourceLimit := &v1alpha1.SwitchResourceLimit{}
		switchResourceLimit.Name = "user-limit"
		switchResourceLimit.Namespace = limit.Namespace
		switchResourceLimit.Status.VLANRange = limit.VLANRange
		switchResourceLimit.Status.SwitchResourceRef.Name = i.Name
		switchResourceLimit.Status.SwitchResourceRef.Namespace = i.Namespace

		err = info.Client.Create(ctx, switchResourceLimit)
		if err != nil {
			// If SwitchResourceLimit is existed, skip it
			if errors.IsAlreadyExists(err) {
				continue
			}
			return machine.ResultContinue(v1alpha1.SwitchResourceCreating, requeueAfterTime, err)
		}

		// updates the value to `status.availableVLAN`.
		err = i.Shrink(limit)
		if err != nil {
			return machine.ResultContinue(v1alpha1.SwitchResourceCreating, 0, err)
		}
	}

	return machine.ResultContinue(v1alpha1.SwitchResourceRunning, 0, nil)
}

func (r *SwitchResourceReconciler) runningHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchResource)

	if !i.DeletionTimestamp.IsZero() {
		return machine.ResultContinue(v1alpha1.SwitchResourceDeleting, 0, nil)
	}

	if !reflect.DeepEqual(i.Spec.TenantLimits, i.Status.TenantLimits) {
		return machine.ResultContinue(v1alpha1.SwitchResourceVerifying, 0, nil)
	}

	// Check switchResourceLimit are existed
	for _, limit := range i.Status.TenantLimits {
		// Get switchResourceLimit
		err := info.Client.Get(
			ctx, types.NamespacedName{
				Name:      "user-limit",
				Namespace: limit.Namespace,
			},
			&v1alpha1.SwitchResourceLimit{},
		)
		if err != nil {
			// If switchResourceLimit isn't find, return creating state and create it
			if errors.IsNotFound(err) {
				return machine.ResultContinue(v1alpha1.SwitchResourceCreating, 0, err)
			}
			return machine.ResultContinue(v1alpha1.SwitchResourceRunning, requeueAfterTime, err)
		}
	}

	return machine.ResultContinue(v1alpha1.SwitchResourceRunning, requeueAfterTime, nil)
}

func (r *SwitchResourceReconciler) deletingHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	i := instance.(*v1alpha1.SwitchResource)

	// Foreground delete
	propagationPolicy := metav1.DeletePropagationForeground
	err := info.Client.Delete(ctx, i, &client.DeleteOptions{PropagationPolicy: &propagationPolicy})
	if err != nil {
		return v1alpha1.SwitchResourceDeleting, ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	// Remove finalizer
	finalizer.Remove(&i.Finalizers, finalizerKey)

	return machine.ResultComplete(v1alpha1.SwitchResourceDeleting, nil)
}
