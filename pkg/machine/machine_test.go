package machine

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type testInstance struct {
	out   string
	state StateType
	runtime.Unknown
}

func (t *testInstance) GetState() StateType {
	return t.state
}

func (t *testInstance) SetState(state StateType) {
	t.state = state
}

func (t *testInstance) SetError(err error) {

}

func handlerTest0(ctx context.Context, info *ReconcileInfo, instance interface{}) (StateType, ctrl.Result, error) {
	instance.(*testInstance).out = "Hello"
	return "test1", ctrl.Result{}, nil
}

func handlerTest1(ctx context.Context, info *ReconcileInfo, instance interface{}) (StateType, ctrl.Result, error) {
	instance.(*testInstance).out += " world"
	return "test2", ctrl.Result{}, nil
}

func handlerTest2(ctx context.Context, info *ReconcileInfo, instance interface{}) (StateType, ctrl.Result, error) {
	instance.(*testInstance).out += "!"
	return "", ctrl.Result{}, nil
}

func TestMachine(t *testing.T) {
	var instance testInstance
	m := New(
		nil,
		&instance,
		map[StateType]Handler{
			"":      handlerTest0,
			"test1": handlerTest1,
			"test2": handlerTest2,
		},
	)
	_, _, _ = m.Reconcile(context.TODO())
	if instance.out != "Hello" {
		t.Fatal(instance.out)
	}
	_, _, _ = m.Reconcile(context.TODO())
	if instance.out != "Hello world" {
		t.Fatal(instance.out)
	}
	_, _, _ = m.Reconcile(context.TODO())
	if instance.out != "Hello world!" {
		t.Fatal(instance.out)
	}
}
