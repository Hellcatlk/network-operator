package machine

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

// ResultContinue continue reconcile after delay
func ResultContinue(state StateType, delay time.Duration, err error) (StateType, ctrl.Result, error) {
	return state, ctrl.Result{Requeue: true, RequeueAfter: delay}, err
}

// ResultComplete stop reconcile
func ResultComplete(state StateType, err error) (StateType, ctrl.Result, error) {
	return state, ctrl.Result{}, err
}
