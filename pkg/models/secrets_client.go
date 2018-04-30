package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// SecretsClient is a convenience wrapper for managing k8s Secrets
type SecretsClient struct {
	secretInterface apiv1.SecretInterface
	Namespace       string
	AuthInfo        string
}

// NewSecretsClient constructor
func NewSecretsClient(namespace string) (*SecretsClient, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		var err error
		namespace, _, err = kubeConfig.Namespace()
		if err != nil {
			return nil, err
		}
	}

	return &SecretsClient{
		secretInterface: client.CoreV1().Secrets(namespace),
		Namespace:       namespace,
		AuthInfo:        rawConfig.Contexts[rawConfig.CurrentContext].AuthInfo,
	}, nil
}

// List all Secrets
func (s *SecretsClient) List() (*v1.SecretList, error) {
	return s.secretInterface.List(metav1.ListOptions{})
}

// Create a new Secret
func (s *SecretsClient) Create(name string) (*v1.Secret, error) {
	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return s.secretInterface.Create(&secret)
}

// CreateWithData creates a new Secret and passed in Data keys
func (s *SecretsClient) CreateWithData(name string, data map[string][]byte) (*v1.Secret, error) {

	// TODO: add key annotations

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	return s.secretInterface.Create(&secret)
}

// Delete a secret
func (s *SecretsClient) Delete(name string) error {
	return s.secretInterface.Delete(name, &metav1.DeleteOptions{})
}

// Get Secret
func (s *SecretsClient) Get(name string) (*v1.Secret, error) {
	return s.secretInterface.Get(name, metav1.GetOptions{})
}

// Update Secret keys
func (s *SecretsClient) Update(name string, data map[string][]byte) (*v1.Secret, error) {
	secret, err := s.Get(name)
	if err != nil {
		return s.CreateWithData(name, data)
	}

	// TODO: update key annotations

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data = data

	return s.secretInterface.Update(secret)
}
