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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	metal3iov1alpha1 "github.com/metal3-io/networkconfiguration-operator/api/v1alpha1"
	"github.com/metal3-io/networkconfiguration-operator/pkg/machine"
)

// SwitchPortReconciler reconciles a SwitchPort object
type SwitchPortReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal3.io,resources=switchports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=switchports/status,verbs=get;update;patch

// Reconcile ...
func (r *SwitchPortReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	ctx := context.Background()
	logger := r.Log.WithValues("switchport", req.NamespacedName)

	// Fetch the instance
	instance := &v1alpha1.SwitchPort{}
	err = r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		// Error reading the object - requeue the request
		return reconcile.Result{}, err
	}

	if len(instance.OwnerReferences) == 0 {
		return result, fmt.Errorf("The OwnerReferences of port mustn't be empty")
	}

	// Initialize state machine
	m := machine.New(
		&machine.Information{
			Client: r.Client,
			Logger: logger,
		},
		instance,
		&machine.Handlers{
			v1alpha1.SwitchPortNone:        r.noneHandler,
			v1alpha1.SwitchPortIdle:        r.idleHandler,
			v1alpha1.SwitchPortValidating:  r.validatingandler,
			v1alpha1.SwitchPortConfiguring: r.configuringHandler,
			v1alpha1.SwitchPortActive:      r.activeHandler,
			v1alpha1.SwitchPortCleaning:    r.cleaningHandler,
			v1alpha1.SwitchPortDeleting:    r.deletingHandler,
		},
	)

	// Reconcile state machine
	result, merr := m.Reconcile(ctx)
	if merr != nil {
		err = merr.Error()
		switch merr.Type() {
		case machine.ReconcileError:
			logger.Error(err, "reconcile error")
		case machine.HandlerError:
			logger.Error(err, "handler error")
		}
	}

	// Update object
	err = r.Update(ctx, instance)

	return ctrl.Result{}, err
}

// SetupWithManager ...
func (r *SwitchPortReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3iov1alpha1.SwitchPort{}).
		Complete(r)
}
