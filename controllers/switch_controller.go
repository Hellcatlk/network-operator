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

	metal3iov1alpha1 "github.com/Hellcatlk/network-operator/api/v1alpha1"
	"github.com/Hellcatlk/network-operator/pkg/machine"
)

// SwitchReconciler reconciles a Switch object
type SwitchReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal3.io,resources=switches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=switches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal3.io,resources=switches/finalizers,verbs=update

// +kubebuilder:rbac:groups=metal3.io,resources=ansibleswitches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=ansibleswitches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal3.io,resources=ansibleswitches/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete

// Reconcile switch resources
func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("switch", req.NamespacedName)

	// Fetch the instance
	instance := &metal3iov1alpha1.Switch{}
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
			metal3iov1alpha1.SwitchNone:        r.noneHandler,
			metal3iov1alpha1.SwitchVerifying:   r.verifyingHandler,
			metal3iov1alpha1.SwitchConfiguring: r.configuringHandler,
			metal3iov1alpha1.SwitchRunning:     r.runningHandler,
			metal3iov1alpha1.SwitchDeleting:    r.deletingHandler,
		},
	)

	// Reconcile state machine
	dirty, result, err := m.Reconcile(ctx)
	if err != nil {
		logger.Error(err, "state machine error")
	}

	// Only need to update switch port when it dirty
	if dirty == machine.MetadataAndSpec || dirty == machine.All {
		logger.Info("updating switch")
		err = r.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "update switch failed")
			return result, err
		}
	}
	if dirty == machine.Status || dirty == machine.All {
		logger.Info("updating switch status")
		err = r.Status().Update(ctx, instance)
		if err != nil {
			logger.Error(err, "update switch status failed")
			return result, err
		}
	}

	return result, err
}

// SetupWithManager register reconciler
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3iov1alpha1.Switch{}).
		Complete(r)
}
