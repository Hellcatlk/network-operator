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
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=switchports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=switchports/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal3.io,resources=ovsswitches,verbs=get;list;watch;create;update;patch;delete

// Reconcile switch resources
func (r *SwitchReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("switch", req.NamespacedName)

	// Fetch the instance
	instance := &v1alpha1.Switch{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		// The object has been deleted
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		return ctrl.Result{}, err
	}

	// Initialize state machine
	m := machine.New(
		&machine.ReconcileInfo{
			Client: r.Client,
			Logger: logger,
		},
		instance,
		&machine.Handlers{
			metal3iov1alpha1.SwitchNone:        r.noneHandler,
			metal3iov1alpha1.SwitchVerify:      r.verifyingHandler,
			metal3iov1alpha1.SwitchConfiguring: r.configuringHandler,
			metal3iov1alpha1.SwitchRunning:     r.runningHandler,
			metal3iov1alpha1.SwitchDeleting:    r.deletingHandler,
		},
	)

	// Reconcile state machine
	dirty, result, err := m.Reconcile(ctx)

	// Only update switch port when it dirty
	if dirty {
		logger.Info("updating switch")
		err = r.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "update switch failed")
		}
	}

	return result, nil
}

// SetupWithManager register reconciler
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3iov1alpha1.Switch{}).
		Complete(r)
}
