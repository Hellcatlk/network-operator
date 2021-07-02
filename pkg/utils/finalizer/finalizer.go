package finalizer

import (
	"github.com/Hellcatlk/network-operator/pkg/utils/strings"
)

// Add create finalizer, must unique for every object.
func Add(finalizers *[]string, finalizer string) {
	if strings.SliceContains(*finalizers, finalizer) {
		return
	}

	*finalizers = append(*finalizers, finalizer)
}

// Remove remove finalizer
func Remove(finalizers *[]string, finalizer string) {
	strings.SliceDelete(finalizers, finalizer)
}
