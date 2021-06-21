package certificate

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Certificate include username and password
type Certificate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Fetch secret
func Fetch(ctx context.Context, client client.Client, secretRef *corev1.SecretReference) (*Certificate, error) {
	if secretRef == nil {
		return nil, fmt.Errorf("secret reference is nil")
	}

	instance := &corev1.Secret{}
	err := client.Get(
		ctx,
		types.NamespacedName{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
		},
		instance,
	)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		Username: string(instance.Data["username"]),
		Password: string(instance.Data["password"]),
	}, nil
}
