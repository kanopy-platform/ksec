package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [secret...]",
	Short: "Create a Secret",
	Args:  cobra.MinimumNArgs(1),
	Run:   createCommand,
}

func createCommand(cmd *cobra.Command, args []string) {
	for _, name := range args {
		if _, err := secretsClient.Create(name); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Created secret \"%s\"\n", name)
	}
}
