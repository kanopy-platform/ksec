package models

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var secretsClient *SecretsClient
var expectedSecretName = "test-secret"
var defaultNamespace = "default"

// setupTestClient initializes a new secretsClient before each test
func setupTestClient(namespace string) {
	secretsClient = MockNewSecretsClient(namespace)
}

func TestNewSecretsClient(t *testing.T) {
	setupTestClient(defaultNamespace)
	assert.Equal(t, "default", secretsClient.Namespace, "Namespace should be %s", defaultNamespace)
}

func TestCreate(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	secret, err := secretsClient.Create(ctx, expectedSecretName)
	assert.NoError(t, err, "Creating secret should not return an error")

	assert.NotNil(t, secret, "Created secret should not be nil")
	assert.Equal(t, secret.Name, expectedSecretName)
}

func TestList(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	_, err := secretsClient.Create(ctx, expectedSecretName)
	assert.NoError(t, err, "Creating secret should not return an error")

	secrets, err := secretsClient.List(ctx)
	assert.NoError(t, err, "Listing secrets should not return an error")
	assert.NotEmpty(t, secrets.Items, "Secrets list should not be empty")
	assert.Equal(t, expectedSecretName, secrets.Items[0].Name, "First secret name should be %s", expectedSecretName)
}

func TestCreateWithData(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	dataArgs := []string{
		"key=value",
		"secret-key=secret-value",
		"ENV_VAR=~!@#$%^&*()_+",
		"DB_URL=mongodb://host1.example.com:27017,host2.example.com:27017/prod?replicaSet=prod",
	}

	data := make(map[string][]byte)

	for _, item := range dataArgs {
		split := strings.SplitN(item, "=", 2)
		assert.Len(t, split, 2, "Data is not formatted correctly")
		data[split[0]] = []byte(split[1])
	}

	secret, err := secretsClient.CreateWithData(ctx, expectedSecretName, data)
	assert.NoError(t, err, "Creating secret with data should not return an error")

	assert.NotNil(t, secret.Data, "Created secret with data should not be nil")
	assert.Equal(t, data, secret.Data)
}

func TestGet(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	_, err := secretsClient.CreateWithData(ctx, expectedSecretName, map[string][]byte{
		"DB_URL": []byte("mongodb://host1.example.com:27017,host2.example.com:27017/prod?replicaSet=prod"),
	})
	assert.NoError(t, err, "Creating secret with data should not return an error")

	secret, err := secretsClient.Get(ctx, expectedSecretName)
	assert.NoError(t, err, "Getting secret should not return an error")
	assert.Equal(t, "mongodb://host1.example.com:27017,host2.example.com:27017/prod?replicaSet=prod", string(secret.Data["DB_URL"]), "DB_URL should match the expected value")
}

func TestGetKey(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	_, err := secretsClient.CreateWithData(ctx, expectedSecretName, map[string][]byte{
		"secret-key": []byte("secret-value"),
	})
	assert.NoError(t, err, "Creating secret with data should not return an error")

	value, err := secretsClient.GetKey(ctx, expectedSecretName, "secret-key")
	assert.NoError(t, err, "Getting secret key should not return an error")
	assert.Equal(t, "secret-value", value, "Secret key value should match the expected value")

	value, err = secretsClient.GetKey(ctx, expectedSecretName, "thiskeydoesnotexist")
	assert.Error(t, err, "Getting a non-existent key should return an error")
	assert.Empty(t, value, "Non-existent key should return an empty value")
}

func TestUpdate(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	secret, err := secretsClient.CreateWithData(ctx, expectedSecretName, map[string][]byte{
		"key": []byte("oldvalue"),
	})
	assert.NoError(t, err, "Creating secret with data should not return an error")

	data := map[string][]byte{
		"key": []byte("newvalue"),
	}

	secret, err = secretsClient.Update(ctx, secret, data)
	assert.NoError(t, err, "Updating secret should not return an error")
	assert.Equal(t, "newvalue", string(secret.Data["key"]), "Key value should be updated to 'newvalue'")
}

func TestUpsert(t *testing.T) {
	setupTestClient(defaultNamespace)
	ctx := context.Background()

	data := map[string][]byte{
		"key": []byte("upsert"),
	}

	// First upsert (should create the secret)
	secret, err := secretsClient.Upsert(ctx, expectedSecretName, data)
	assert.NoError(t, err, "Upserting (creating) secret should not return an error")
	assert.Equal(t, "upsert", string(secret.Data["key"]), "Key value should be 'upsert'")

	// Second upsert (should update the secret)
	data["key"] = []byte("upserted")
	secret, err = secretsClient.Upsert(ctx, expectedSecretName, data)
	assert.NoError(t, err, "Upserting (updating) secret should not return an error")
	assert.Equal(t, "upserted", string(secret.Data["key"]), "Key value should be 'upserted'")
}
