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

// PortReconciler reconciles a Port object
type PortReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal3.io,resources=ports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=ports/status,verbs=get;update;patch

// Reconcile ...
func (r *PortReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	_ = context.Background()
	_ = r.Log.WithValues("port", req.NamespacedName)

	_ = context.Background()
	_ = r.Log.WithValues("Port", req.NamespacedName)

	// Fetch the instance
	instance := &v1alpha1.Port{}
	err = r.Get(context.TODO(), req.NamespacedName, instance)
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
			Logger: r.Log,
		},
		instance,
		&machine.Handlers{
			v1alpha1.PortNone:        r.noneHandler,
			v1alpha1.PortCreated:     r.createdHandler,
			v1alpha1.PortConfiguring: r.configuringHandler,
			v1alpha1.PortConfigured:  r.configuredHandler,
			v1alpha1.PortCleaning:    r.cleaningHandler,
			v1alpha1.PortCleaned:     r.cleanedHandler,
		},
	)

	var merr *machine.Error
	switch {
	// On object created
	case instance.DeletionTimestamp.IsZero() && len(instance.Finalizers) == 0:
		// Reconcile state
		result, merr = m.Reconcile(context.TODO())
		if merr != nil {
			err = merr.Error()
			switch merr.Type() {
			case machine.ReconcileError:
				// Do something
			case machine.HandlerError:
				// Do something
			}
		}

	// On object updated
	case instance.DeletionTimestamp.IsZero():
		// Reconcile state
		result, merr = m.Reconcile(context.TODO())
		if merr != nil {
			err = merr.Error()
			switch merr.Type() {
			case machine.ReconcileError:
				// Do something
			case machine.HandlerError:
				// Do something
			}
		}

	// On object delete
	case !instance.DeletionTimestamp.IsZero():
		// Set instance's state to `Cleaning`
		instance.SetState(v1alpha1.PortCleaning)
		// Reconcile state
		result, merr = m.Reconcile(context.TODO())
		if merr != nil {
			err = merr.Error()
			switch merr.Type() {
			case machine.ReconcileError:
				// Do something
			case machine.HandlerError:
				// Do something
			}
		}
	}

	// Update object
	err = r.Update(context.TODO(), instance)

	return ctrl.Result{}, err
}

// SetupWithManager ...
func (r *PortReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3iov1alpha1.Port{}).
		Complete(r)
}
