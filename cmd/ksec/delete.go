package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret...]",
	Short: "Delete a Secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  deleteCommand,
}

func deleteCommand(cmd *cobra.Command, args []string) error {
	skipconfirm, err := cmd.Flags().GetBool("yes")
	if err != nil {
		return err
	}

	for _, name := range args {
		if _, err := secretsClient.Get(name); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		confirmationMessage := fmt.Sprintf(`Delete secret "%s"? This action cannot be reversed.`, name)
		if !skipconfirm && !askConfirmation(confirmationMessage) {
			fmt.Println("Delete canceled")
			continue
		}

		if err := secretsClient.Delete(name); err != nil {
			return err
		}
		fmt.Printf("Deleted secret \"%s\"\n", name)
	}
	return nil
}
