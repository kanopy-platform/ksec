package models

import (
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func MockNewSecretsClient(config clientcmd.ClientConfig, namespace string) (*SecretsClient, error) {
	if namespace == "" {
		_, _, err := config.Namespace()
		if err != nil {
			return nil, err
		}
	}

	return &SecretsClient{
		secretInterface: testclient.NewSimpleClientset().CoreV1().Secrets(namespace),
		Namespace:       namespace,
		AuthInfo:        "testuser",
	}, nil
}

// Mock ClientConfig for context checks
func MockClientConfig() clientcmd.ClientConfig {
	return clientcmd.NewDefaultClientConfig(api.Config{
		Contexts: map[string]*api.Context{
			"default-context": {
				Namespace: "default",
			},
		},
		CurrentContext: "default-context",
	}, &clientcmd.ConfigOverrides{})
}
