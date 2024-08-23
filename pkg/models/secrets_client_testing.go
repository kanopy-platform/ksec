package models

import (
	testclient "k8s.io/client-go/kubernetes/fake"
)

// MockNewSecretsClient creates a new SecretsClient with a mock Kubernetes client
func MockNewSecretsClient(namespace string) *SecretsClient {
	return &SecretsClient{
		secretInterface: testclient.NewSimpleClientset().CoreV1().Secrets(namespace),
		Namespace:       namespace,
		AuthInfo:        "testuser",
	}
}
