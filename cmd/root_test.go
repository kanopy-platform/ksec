package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/colinhoglund/ksec/pkg/models"
	"github.com/spf13/cobra"
)

// mock rootCmd
var rootCmd *cobra.Command

func TestMain(m *testing.M) {
	rootCmd = &cobra.Command{
		Use:     "ksec",
		Short:   "A tool for managing Kubernetes Secret data",
		Version: "0.1.0",
	}
	initRootCmd(rootCmd)
	secretsClient = models.MockNewSecretsClient()
	m.Run()
}

//helpers
func testErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func cmdExec(args []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		return &bytes.Buffer{}, err
	}
	return buf, nil
}

// tests
func TestCreateSecret(t *testing.T) {
	_, err := cmdExec([]string{"create", "test"})
	testErr(err, t)

	_, err = cmdExec([]string{"set", "test", "key=value"})
	testErr(err, t)

	secret, err := secretsClient.Get("test")
	testErr(err, t)

	val, ok := secret.Data["key"]
	if !ok {
		t.Fatal("Key does not exist")
	} else if string(val) != "value" {
		t.Fatal("Key has incorrect value")
	}
}

func TestDeleteSecret(t *testing.T) {
	_, err := cmdExec([]string{"delete", "test"})
	testErr(err, t)

	_, err = cmdExec([]string{"get", "test"})
	if !strings.HasSuffix(err.Error(), "not found") {
		t.Fatal("Secret still exists")
	}
}
