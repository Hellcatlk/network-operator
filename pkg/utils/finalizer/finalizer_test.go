package finalizer

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type testType struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func TestAdd(t *testing.T) {
	var object testType

	cases := []struct {
		finalizer string
		expected  []string
	}{
		{
			finalizer: "test1",
			expected:  []string{"test1"},
		},
		{
			finalizer: "test2",
			expected:  []string{"test1", "test2"},
		},
		{
			finalizer: "test2",
			expected:  []string{"test1", "test2"},
		},
	}

	for _, c := range cases {
		t.Run(c.finalizer, func(t *testing.T) {
			Add(&object.Finalizers, c.finalizer)
			if !reflect.DeepEqual(c.expected, object.Finalizers) {
				t.Errorf("expected: %v, got: %v", c.expected, object.Finalizers)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	var object testType
	object.Finalizers = []string{"test1", "test2", "test3"}

	cases := []struct {
		finalizer string
		expected  []string
	}{
		{
			finalizer: "test1",
			expected:  []string{"test2", "test3"},
		},
		{
			finalizer: "test2",
			expected:  []string{"test3"},
		},
		{
			finalizer: "test2",
			expected:  []string{"test3"},
		},
	}

	for _, c := range cases {
		t.Run(c.finalizer, func(t *testing.T) {
			Remove(&object.Finalizers, c.finalizer)
			if !reflect.DeepEqual(c.expected, object.Finalizers) {
				t.Errorf("expected: %v, got: %v", c.expected, object.Finalizers)
			}
		})
	}
}
