package models

import (
	testclient "k8s.io/client-go/kubernetes/fake"
)

func MockNewSecretsClient(namespace string) *SecretsClient {
	return &SecretsClient{
		secretInterface: testclient.NewSimpleClientset().CoreV1().Secrets(namespace),
		Namespace:       namespace,
		AuthInfo:        "testuser",
	}
}
