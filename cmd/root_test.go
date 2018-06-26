package cmd

import (
	"testing"

	"github.com/colinhoglund/ksec/pkg/models"
	"github.com/spf13/cobra"
)

func MockNewRootCmd() *cobra.Command {
	rootCmd = &cobra.Command{
		Use:     "ksec",
		Short:   "A tool for managing Kubernetes Secret data",
		Version: "0.1.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			secretsClient = models.MockNewSecretsClient()
		},
	}

	initRootCmd()
	return rootCmd
}

func TestRootCmd(t *testing.T) {
	if err := MockNewRootCmd().Execute(); err != nil {
		t.Fatal(err.Error())
	}
}
