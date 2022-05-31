package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all secrets in a namespace",
	Args:    cobra.NoArgs,
	RunE:    listCommand,
}

func listCommand(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	secrets, err := secretsClient.List(ctx)
	if err != nil {
		return err
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
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
			if secret.Type == v1.SecretTypeOpaque {
				fmt.Println(secret.Name)
			}
		}
	}
	return nil
}
