package models

import (
	"encoding/json"
	"fmt"

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
	// initialize secrets client
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace, _, err = kubeConfig.Namespace()
		if err != nil {
			return nil, err
		}
	}

	return &SecretsClient{
		secretInterface: clientSet.CoreV1().Secrets(namespace),
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
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return s.secretInterface.Create(&secret)
}

// CreateWithData creates a new Secret and passed in Data keys
func (s *SecretsClient) CreateWithData(name string, data map[string][]byte) (*v1.Secret, error) {

	annotation := NewKeyAnnotation(s.AuthInfo)
	annotations := make(map[string]string)

	for key := range data {
		jsonBytes, err := json.Marshal(annotation)
		if err != nil {
			return nil, err
		}
		annotations[fmt.Sprintf("%s/%s", annotationPrefix, key)] = string(jsonBytes)
	}

	secret := v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
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

// GetKey retrieves an individual keys value from a secret
func (s *SecretsClient) GetKey(name, key string) (string, error) {
	secret, err := s.Get(name)
	if err != nil {
		return "", err
	}

	value, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("secret key %s does not exist", key)
	}

	return string(value), nil
}

// Update Secret keys
func (s *SecretsClient) Update(secret *v1.Secret, data map[string][]byte) (*v1.Secret, error) {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}

	annotation := NewKeyAnnotation(s.AuthInfo)
	for key, value := range data {
		secret.Data[key] = value
		jsonBytes, err := json.Marshal(annotation)
		if err != nil {
			return nil, err
		}
		secret.Annotations[fmt.Sprintf("%s/%s", annotationPrefix, key)] = string(jsonBytes)
	}

	return s.secretInterface.Update(secret)
}

// Upsert creates a Secret if needed and updates Secret keys
func (s *SecretsClient) Upsert(name string, data map[string][]byte) (*v1.Secret, error) {
	secret, err := s.Get(name)
	if err != nil {
		return s.CreateWithData(name, data)
	}
	return s.Update(secret, data)
}
