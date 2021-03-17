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

func TestAddFinalizers(t *testing.T) {
	var object testType

	cases := []struct {
		finalizer     string
		expected      []string
		expectedError bool
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
			finalizer:     "test2",
			expected:      []string{"test1", "test2"},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.finalizer, func(t *testing.T) {
			err := AddFinalizer(&object.Finalizers, c.finalizer)
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
			if !reflect.DeepEqual(c.expected, object.Finalizers) {
				t.Errorf("expected: %v, got: %v", c.expected, object.Finalizers)
			}
		})
	}
}

func TestRemoveFinalizers(t *testing.T) {
	var object testType
	object.Finalizers = []string{"test1", "test2", "test3"}

	cases := []struct {
		finalizer     string
		expected      []string
		expectedError bool
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
			finalizer:     "test2",
			expected:      []string{"test3"},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.finalizer, func(t *testing.T) {
			err := RemoveFinalizer(&object.Finalizers, c.finalizer)
			if (err != nil) != c.expectedError {
				t.Errorf("got unexpected error: %v", err)
			}
			if !reflect.DeepEqual(c.expected, object.Finalizers) {
				t.Errorf("expected: %v, got: %v", c.expected, object.Finalizers)
			}
		})
	}
}
