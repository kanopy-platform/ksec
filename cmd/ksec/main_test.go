package main

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/kanopy-platform/ksec/pkg/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// mock rootCmd
var rootCmd *cobra.Command

func TestMain(m *testing.M) {
	rootCmd = &cobra.Command{
		Use:   "ksec",
		Short: "A tool for managing Kubernetes Secret data",
	}
	initRootCmd(rootCmd)
	mockConfig := models.MockClientConfig()
	secretsClient, _ = models.MockNewSecretsClient(mockConfig, "default")
	os.Exit(m.Run())
}

func cmdExec(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

// tests
func TestCreateSecret(t *testing.T) {
	ctx := context.Background()

	err := cmdExec([]string{"create", "test"})
	assert.NoError(t, err, "Creating secret should not return an error")

	err = cmdExec([]string{"set", "test", "key=value"})
	assert.NoError(t, err, "Setting secret key should not return an error")

	secret, err := secretsClient.Get(ctx, "test")
	assert.NoError(t, err, "Getting secret should not return an error")
	assert.NotNil(t, secret, "Secret should not be nil")

	val, ok := secret.Data["key"]
	assert.True(t, ok, "Key should exist in secret data")
	assert.Equal(t, "value", string(val), "Key should have the correct value")
}

func TestUnsetSecretKey(t *testing.T) {
	ctx := context.Background()

	err := cmdExec([]string{"unset", "test", "key"})
	assert.NoError(t, err, "Unsetting secret key should not return an error")

	secret, err := secretsClient.Get(ctx, "test")
	assert.NoError(t, err, "Getting secret should not return an error")
	assert.NotNil(t, secret, "Secret should not be nil")

	_, ok := secret.Data["key"]
	assert.False(t, ok, "Key should not exist in secret data")
}

func TestDeleteSecret(t *testing.T) {
	err := cmdExec([]string{"delete", "test", "--yes"})
	assert.NoError(t, err, "Deleting secret should not return an error")

	err = cmdExec([]string{"get", "test"})
	assert.Error(t, err, "Getting deleted secret should return an error")
	assert.True(t, strings.HasSuffix(err.Error(), "not found"), "Error should indicate that the secret was not found")
}

func TestPushSecret(t *testing.T) {
	ctx := context.Background()
	content := []byte("ENV_VAR=secret")

	tempfile, err := os.CreateTemp("", "ksec")
	assert.NoError(t, err, "Creating temp file should not return an error")
	defer os.Remove(tempfile.Name())

	_, err = tempfile.Write(content)
	assert.NoError(t, err, "Writing to temp file should not return an error")

	err = cmdExec([]string{"push", tempfile.Name(), "pushtest"})
	assert.NoError(t, err, "Pushing secret should not return an error")

	secret, err := secretsClient.Get(ctx, "pushtest")
	assert.NoError(t, err, "Getting pushed secret should not return an error")
	assert.NotNil(t, secret, "Pushed secret should not be nil")

	val, ok := secret.Data["ENV_VAR"]
	assert.True(t, ok, "ENV_VAR key should exist in pushed secret data")
	assert.Equal(t, "secret", string(val), "ENV_VAR key should have the correct value")

	err = tempfile.Close()
	assert.NoError(t, err, "Closing temp file should not return an error")
}

func TestPullSecret(t *testing.T) {
	tempfile, err := os.CreateTemp("", "ksec")
	assert.NoError(t, err, "Creating temp file should not return an error")
	defer os.Remove(tempfile.Name())

	err = cmdExec([]string{"set", "pulltest", "ENV_VAR=secret"})
	assert.NoError(t, err, "Setting secret should not return an error")

	err = cmdExec([]string{"pull", "pulltest", tempfile.Name()})
	assert.NoError(t, err, "Pulling secret should not return an error")

	file, err := os.Open(tempfile.Name())
	assert.NoError(t, err, "Opening temp file should not return an error")
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	assert.NoError(t, err, "Reading line from temp file should not return an error")

	assert.Equal(t, "ENV_VAR=secret", string(line), "File should contain the pulled secret contents")
}
