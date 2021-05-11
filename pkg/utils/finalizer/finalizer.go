package finalizer

import (
	"github.com/metal3-io/networkconfiguration-operator/pkg/utils/stringslice"
)

// Add create finalizer, must unique for every object.
func Add(finalizers *[]string, finalizer string) {
	if stringslice.Contains(*finalizers, finalizer) {
		return
	}

	*finalizers = append(*finalizers, finalizer)
}

// Remove remove finalizer
func Remove(finalizers *[]string, finalizer string) {
	stringslice.Delete(finalizers, finalizer)
}
