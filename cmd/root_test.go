package cmd

import (
	"bytes"
	"testing"

	"github.com/colinhoglund/ksec/pkg/models"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// mock client
func MockNewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ksec",
		Short:   "A tool for managing Kubernetes Secret data",
		Version: "0.1.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			secretsClient = models.MockNewSecretsClient()
		},
	}

	initRootCmd(rootCmd)
	return rootCmd
}

//helpers
func testErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func exec(args []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		return &bytes.Buffer{}, err
	}
	return buf, nil
}

// tests
func TestRootCmd(t *testing.T) {
	rootCmd = MockNewRootCmd()
	_, err := exec([]string{"list"})
	testErr(err, t)
}

func TestCreateCmd(t *testing.T) {
	_, err := exec([]string{"create", "test"})
	testErr(err, t)
}
