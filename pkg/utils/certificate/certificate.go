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
func Fetch(ctx context.Context, client client.Client, ref *corev1.SecretReference) (cert *Certificate, err error) {
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
	if err != nil {
		return nil, err
	}

	cert = &Certificate{}

	cert.Username, err = decode(instance, "username")
	if err != nil {
		return nil, err
	}

	cert.Password, err = decode(instance, "password")
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// decode key's value and decode it by base64
func decode(secret *corev1.Secret, key string) (string, error) {
	if secret == nil {
		return "", fmt.Errorf("secret is nil")
	}

	// Decode to byte
	bytes := make([]byte, base64.StdEncoding.DecodedLen(len(secret.Data[key])))
	len, err := base64.StdEncoding.Decode(bytes, secret.Data[key])
	if err != nil {
		return "", err
	}

	return string(bytes[:len]), nil
}
