package models

import (
	testclient "k8s.io/client-go/kubernetes/fake"
)

// mock client
func MockNewSecretsClient() *SecretsClient {
	return &SecretsClient{
		secretInterface: testclient.NewSimpleClientset().CoreV1().Secrets("default"),
		Namespace:       "default",
		AuthInfo:        "testuser",
	}
}
