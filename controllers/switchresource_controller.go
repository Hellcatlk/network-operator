/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/machine"
)

// SwitchResourceReconciler reconciles a SwitchResource object
type SwitchResourceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal3.io,resources=switchresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=switchresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal3.io,resources=switchresources/finalizers,verbs=update

// Reconcile switch resources
func (r *SwitchResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("switchresource", req.NamespacedName)

	// Fetch the instance
	instance := &v1alpha1.SwitchResource{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		// The object has been deleted
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Requeue when other error
		return ctrl.Result{}, err
	}

	// Initialize state machine
	m := machine.New(
		&machine.ReconcileInfo{
			Client: r.Client,
			Logger: logger,
		},
		instance,
		map[machine.StateType]machine.Handler{
			v1alpha1.SwitchResourceNone:      r.noneHandler,
			v1alpha1.SwitchResourceVerifying: r.verifyingHandler,
			v1alpha1.SwitchResourceCreating:  r.creatingHandler,
			v1alpha1.SwitchResourceRunning:   r.runningHandler,
			v1alpha1.SwitchResourceDeleting:  r.deletingHandler,
		},
	)

	// Reconcile state machine
	dirty, result, err := m.Reconcile(ctx)
	if err != nil {
		logger.Error(err, "state machine error")
	}

	// Only need to update switchResource when it dirty
	if dirty == machine.MetadataAndSpec || dirty == machine.All {
		logger.Info("updating switchResource")
		err = r.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "update switchResource failed")
			return result, err
		}
	}
	if dirty == machine.Status || dirty == machine.All {
		logger.Info("updating switchResource status")
		err = r.Status().Update(ctx, instance)
		if err != nil {
			logger.Error(err, "update switchResource status failed")
			return result, err
		}
	}

	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SwitchResource{}).
		Complete(r)
}
