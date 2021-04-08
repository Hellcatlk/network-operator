package finalizer

import (
	"errors"

	"github.com/metal3-io/networkconfiguration-operator/pkg/util/stringslice"
)

// AddFinalizer create finalizer, must unique for every object.
func Add(finalizers *[]string, finalizer string) (err error) {
	if stringslice.Contains(*finalizers, finalizer) {
		return errors.New("the finalizer of object must be unique")
	}
	*finalizers = append(*finalizers, finalizer)
	return
}

// RemoveFinalizer remove finalizer
func Remove(finalizers *[]string, finalizer string) (err error) {
	if !stringslice.Delete(finalizers, finalizer) {
		return errors.New("haven't find finalizer in finalizers")
	}
	return
}
