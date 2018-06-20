package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all secrets in a namespace",
	Args:    cobra.NoArgs,
	Run:     listCommand,
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "Show all secrets (Default: Opaque only)")
}

func listCommand(cmd *cobra.Command, args []string) {
	secrets, err := secretsClient.List()
	if err != nil {
		log.Fatal(err.Error())
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		log.Fatal(err.Error())
	}
	if all {
		lines := []string{"NAME\tTYPE"}
		for _, secret := range secrets.Items {
			lines = append(lines, fmt.Sprintf("%s\t%s", secret.Name, secret.Type))
		}
		outputTabular(lines)
	} else {
		fmt.Println("NAME")
		for _, secret := range secrets.Items {
			if secret.Type == "Opaque" {
				fmt.Println(secret.Name)
			}
		}
	}
}
