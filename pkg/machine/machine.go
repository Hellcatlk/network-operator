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
	handlers map[StateType]Handler
}

// ReconcileInfo is the information need by reconcile
type ReconcileInfo struct {
	Client client.Client
	Logger logr.Logger
}

// Instance is a object for the CR need be reconcile
// NOTE: Instance must be a pointer
type Instance interface {
	runtime.Object
	GetMetadataAndSpec() interface{}
	GetStatus() interface{}
	GetState() StateType
	SetState(state StateType)
	SetError(err error)
}

// Handler is a state handle function. If an error isn't nil or
// ctrl.Result.Requeue is true the state machine will requeue the Request again
type Handler func(ctx context.Context, info *ReconcileInfo, instance interface{}) (StateType, ctrl.Result, error)

// New a state machine
// NOTE: The paramater of instance must be a pointer
func New(info *ReconcileInfo, instance Instance, handlers map[StateType]Handler) Machine {
	return Machine{
		info:     info,
		instance: instance,
		handlers: handlers,
	}
}

// DirtyType means which parts of instance have been modified.
type DirtyType int

const (
	// None means no field changed
	None DirtyType = 0

	// MetadataAndSpec means metadata and spec changed
	MetadataAndSpec DirtyType = 1

	// Status means status
	Status DirtyType = 2

	// All means metadata, spec and status changed
	All DirtyType = 3
)

// Reconcile state machine. If dirty is true, it means the instance has changed
func (m *Machine) Reconcile(ctx context.Context) (DirtyType, ctrl.Result, error) {
	m.info.Logger.Info(string(m.instance.GetState()))

	// There are any handler in handlers?
	if len(m.handlers) == 0 {
		return None, ctrl.Result{}, fmt.Errorf("haven't any handler")
	}

	// Check the state's handler exist or not
	handler, exist := m.handlers[m.instance.GetState()]
	if !exist {
		return None, ctrl.Result{}, fmt.Errorf("no handler for the state(%s)", m.instance.GetState())
	}

	// Call handler
	instanceDeepCopy := m.instance.DeepCopyObject().(Instance)
	nextState, result, err := handler(ctx, m.info, m.instance)
	if err != nil {
		err = fmt.Errorf("%s state handler error: %s", m.instance.GetState(), err)
	}
	m.instance.SetState(nextState)
	m.instance.SetError(err)

	dirty := None
	if !reflect.DeepEqual(m.instance.GetMetadataAndSpec(), instanceDeepCopy.GetMetadataAndSpec()) {
		dirty += MetadataAndSpec
	}
	if !reflect.DeepEqual(m.instance.GetStatus(), instanceDeepCopy.GetStatus()) {
		dirty += Status
	}

	// Check instance is dirty or not
	return dirty, result, nil
}
