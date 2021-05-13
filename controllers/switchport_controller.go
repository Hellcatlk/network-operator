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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Hellcatlk/networkconfiguration-operator/api/v1alpha1"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/machine"
)

const defaultRequeueAfterTime time.Duration = time.Second * 5

// SwitchPortReconciler reconciles a SwitchPort object
type SwitchPortReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal3.io,resources=ovsswitches,verbs=get
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get
// +kubebuilder:rbac:groups=metal3.io,resources=switches,verbs=get
// +kubebuilder:rbac:groups=metal3.io,resources=switchportconfigurations,verbs=get
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
		// The object has been deleted
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		return ctrl.Result{}, err
	}

	if len(instance.OwnerReferences) == 0 || instance.OwnerReferences[0].Kind != "Switch" {
		return result, fmt.Errorf("the OwnerReferences[0] must exist, and it's must be \"Switch\"")
	}

	// Initialize state machine
	m := machine.New(
		&machine.ReconcileInfo{
			Client: r.Client,
			Logger: logger,
		},
		instance,
		&machine.Handlers{
			v1alpha1.SwitchPortNone:        r.noneHandler,
			v1alpha1.SwitchPortIdle:        r.idleHandler,
			v1alpha1.SwitchPortConfiguring: r.configuringHandler,
			v1alpha1.SwitchPortActive:      r.activeHandler,
			v1alpha1.SwitchPortCleaning:    r.cleaningHandler,
			v1alpha1.SwitchPortDeleting:    r.deletingHandler,
		},
	)

	// Reconcile state machine
	dirty, result, merr := m.Reconcile(ctx)
	var errorMessage string
	if merr != nil {
		logger.Error(merr.Error(), string(merr.Type()))
		errorMessage = fmt.Sprintf("%s: %s", merr.Type(), merr.Error())
	}

	// Only update switch port when it dirty
	if dirty || instance.Status.Error != errorMessage {
		logger.Info("updating switch port")
		instance.Status.Error = errorMessage
		err = r.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "update switch port failed")
		}
	}

	// Set default requeue after time
	if result.Requeue && result.RequeueAfter == 0 {
		result.RequeueAfter = defaultRequeueAfterTime
	}
	return result, err
}

// SetupWithManager ...
func (r *SwitchPortReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SwitchPort{}).
		Complete(r)
}
