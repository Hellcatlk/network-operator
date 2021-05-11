package secret

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Fetch secret
func Fetch(ctx context.Context, client client.Client, ref *corev1.SecretReference) (instance *corev1.Secret, err error) {
	if ref == nil {
		return nil, fmt.Errorf("reference is nil")
	}

	instance = &corev1.Secret{}
	err = client.Get(
		ctx,
		types.NamespacedName{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
		instance,
	)

	return
}
