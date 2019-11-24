package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/colinhoglund/ksec/pkg/models"
	"github.com/colinhoglund/ksec/pkg/version"
	"github.com/spf13/cobra"
)

// mock rootCmd
var rootCmd *cobra.Command

func TestMain(m *testing.M) {
	rootCmd = &cobra.Command{
		Use:     "ksec",
		Short:   "A tool for managing Kubernetes Secret data",
		Version: version.Version,
	}
	initRootCmd(rootCmd)
	secretsClient = models.MockNewSecretsClient()
	os.Exit(m.Run())

}

//helpers
func testErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func cmdExec(args []string) error {
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// tests
func TestCreateSecret(t *testing.T) {
	err := cmdExec([]string{"create", "test"})
	testErr(err, t)

	err = cmdExec([]string{"set", "test", "key=value"})
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

func TestUnsetSecretKey(t *testing.T) {
	err := cmdExec([]string{"unset", "test", "key"})
	testErr(err, t)

	secret, err := secretsClient.Get("test")
	testErr(err, t)

	if _, ok := secret.Data["key"]; ok {
		t.Fatal("Key should not exist")
	}
}

func TestDeleteSecret(t *testing.T) {
	err := cmdExec([]string{"delete", "test", "--yes"})
	testErr(err, t)

	err = cmdExec([]string{"get", "test"})
	if !strings.HasSuffix(err.Error(), "not found") {
		t.Fatal("Secret still exists")
	}
}

func TestPushSecret(t *testing.T) {
	content := []byte("ENV_VAR=secret")
	tempfile, err := ioutil.TempFile("", "ksec")
	testErr(err, t)
	defer os.Remove(tempfile.Name())

	_, err = tempfile.Write(content)
	testErr(err, t)

	err = cmdExec([]string{"push", tempfile.Name(), "pushtest"})
	testErr(err, t)

	secret, err := secretsClient.Get("pushtest")
	testErr(err, t)

	val, ok := secret.Data["ENV_VAR"]
	if !ok {
		t.Fatal("Key does not exist")
	} else if string(val) != "secret" {
		t.Fatal("Key has incorrect value")
	}

	err = tempfile.Close()
	testErr(err, t)
}

func TestPullSecret(t *testing.T) {
	tempfile, err := ioutil.TempFile("", "ksec")
	testErr(err, t)
	defer os.Remove(tempfile.Name())

	err = cmdExec([]string{"set", "pulltest", "ENV_VAR=secret"})
	testErr(err, t)

	err = cmdExec([]string{"pull", "pulltest", tempfile.Name()})
	testErr(err, t)

	reader := bufio.NewReader(tempfile)
	line, _, err := reader.ReadLine()
	testErr(err, t)

	if string(line) != "ENV_VAR=secret" {
		t.Fatal("File does not contain pulled contents")
	}
}
