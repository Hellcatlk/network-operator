package machine

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StateType is the type of .status.state
type StateType string

// Machine is a state machine
type Machine struct {
	info     *ReconcileInfo
	instance Instance
	handlers *Handlers
}

// ReconcileInfo is the information need by reconcile
type ReconcileInfo struct {
	Client client.Client
	Logger logr.Logger
}

// Instance is a object for the CR need be reconcile
// NOTE: Instance must be a pointer
type Instance interface {
	GetState() StateType
	SetState(state StateType)
	runtime.Object
}

// Handlers includes a lot of handler
type Handlers map[StateType]Handler

// Handler is a state handle function
type Handler func(ctx context.Context, info *ReconcileInfo, instance interface{}) (nextState StateType, result ctrl.Result, err error)

// ErrorType is the error when reconcile state machine
type ErrorType string

const (
	// ReconcileError means have error when reconcile
	ReconcileError ErrorType = "reconcile error"

	// HandlerError means have error in the handler for a state
	HandlerError ErrorType = "handler error"
)

// Error include error type and error message from state machine
type Error struct {
	errType ErrorType
	err     error
}

// Type return error's type
func (e *Error) Type() ErrorType {
	return e.errType
}

// Error return error itself
func (e *Error) Error() error {
	return e.err
}

// New a state machine
// NOTE: The paramater of instance must be a pointer
func New(info *ReconcileInfo, instance Instance, handlers *Handlers) Machine {
	return Machine{
		info:     info,
		instance: instance,
		handlers: handlers,
	}
}

// Reconcile state machine. If dirty is true, it means the instance has changed,
func (m *Machine) Reconcile(ctx context.Context) (dirty bool, result ctrl.Result, merr *Error) {
	m.info.Logger.Info("reconcile from status %v", m.instance.GetState())

	// Deal possible panic
	defer func() {
		err := recover()
		if err != nil {
			merr = &Error{
				errType: HandlerError,
				err:     fmt.Errorf("handler panic: %v", err),
			}
		}
	}()

	result = ctrl.Result{
		Requeue: false,
	}

	// There are any handler in handlers?
	if m.handlers == nil {
		return dirty, result, &Error{
			errType: ReconcileError,
			err:     fmt.Errorf("haven't any handler"),
		}
	}

	// Check the state's handler exist or not
	handler, exist := (*m.handlers)[m.instance.GetState()]
	if !exist {
		return dirty, result, &Error{
			errType: ReconcileError,
			err:     fmt.Errorf("no handler for the state(%v)", m.instance.GetState()),
		}
	}

	// Call handler
	instanceDeepCopy := m.instance.DeepCopyObject()
	nextState, result, err := handler(ctx, m.info, m.instance)
	m.instance.SetState(nextState)
	if err != nil {
		merr = &Error{
			errType: HandlerError,
			err:     err,
		}
	}

	// Check instance need update or not
	if !reflect.DeepEqual(m.instance, instanceDeepCopy) {
		dirty = true
	}

	return dirty, result, merr
}
