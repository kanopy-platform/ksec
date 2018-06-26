package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [secret...]",
	Short: "Create a Secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  createCommand,
}

func createCommand(cmd *cobra.Command, args []string) error {
	for _, name := range args {
		if _, err := secretsClient.Create(name); err != nil {
			return err
		}
		fmt.Printf("Created secret \"%s\"\n", name)
	}
	return nil
}
