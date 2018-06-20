package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret...]",
	Short: "Delete a Secret",
	Args:  cobra.MinimumNArgs(1),
	Run:   deleteCommand,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func deleteCommand(cmd *cobra.Command, args []string) {
	for _, name := range args {
		if err := secretsClient.Delete(name); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Deleted secret \"%s\"\n", name)
	}
}
