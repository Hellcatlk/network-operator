package certificate

import (
	"context"
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Certificate include username and password
type Certificate struct {
	Username string
	Password string
}

// Fetch secret
func Fetch(ctx context.Context, client client.Client, ref *corev1.SecretReference) (sercet *Certificate, err error) {
	if ref == nil {
		return nil, fmt.Errorf("reference is nil")
	}

	instance := &corev1.Secret{}
	err = client.Get(
		ctx,
		types.NamespacedName{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
		instance,
	)

	return &Certificate{
		Username: decode(instance, "username"),
		Password: decode(instance, "password"),
	}, nil
}

// Parse key and return decoding it by base64
func decode(secret *corev1.Secret, key string) string {
	if secret == nil {
		return ""
	}

	// Decode to byte
	bytes := make([]byte, base64.StdEncoding.DecodedLen(len(secret.Data[key])))
	len, err := base64.StdEncoding.Decode(bytes, []byte(secret.Data[key]))
	if err != nil {
		return ""
	}
	bytes = bytes[:len]

	return string(bytes)
}
