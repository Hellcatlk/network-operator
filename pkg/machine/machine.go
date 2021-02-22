package machine

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StateType is the type of .status.state
type StateType string

// Handler is a state handle function
type Handler func(ctx context.Context, info *Information, instance interface{}) (nextState StateType, result ctrl.Result, err error)

// Handlers includes a lot of handler
type Handlers map[StateType]Handler

// Instance is a object for the CR need be reconcile
// NOTE: Instance must be a pointer
type Instance interface {
	GetState() StateType
	SetState(state StateType)
}

// Information ...
type Information struct {
	Client client.Client
	Logger logr.Logger
}

// Machine is a state machine
type Machine struct {
	info         *Information
	instance     Instance
	handlers     *Handlers
	requeueAfter time.Duration
}

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
func New(info *Information, instance Instance, handlers *Handlers) Machine {
	return Machine{
		info:     info,
		instance: instance,
		handlers: handlers,
	}
}

// Reconcile state machine
func (m *Machine) Reconcile(ctx context.Context) (result ctrl.Result, merr *Error) {
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

	// There are any handler in handlers?
	if m.handlers == nil {
		return result, &Error{
			errType: ReconcileError,
			err:     fmt.Errorf("haven't any handler"),
		}
	}

	// Check the state's handler exist or not
	handler, exist := (*m.handlers)[m.instance.GetState()]
	if !exist {
		return result, &Error{
			errType: ReconcileError,
			err:     fmt.Errorf("no handler for the state(%v)", m.instance.GetState()),
		}
	}

	// Call handler
	nextState, result, err := handler(ctx, m.info, m.instance)
	m.instance.SetState(nextState)
	if err != nil {
		merr = &Error{
			errType: HandlerError,
			err:     err,
		}
	}

	return result, merr
}
