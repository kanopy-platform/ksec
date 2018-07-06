package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret...]",
	Short: "Delete a Secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  deleteCommand,
}

func deleteCommand(cmd *cobra.Command, args []string) error {
	for _, name := range args {
		if err := secretsClient.Delete(name); err != nil {
			return err
		}
		fmt.Printf("Deleted secret \"%s\"\n", name)
	}
	return nil
}
