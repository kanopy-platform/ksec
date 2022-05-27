package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset [secret] [key...]",
	Short: "Unset values in a Secret",
	Args:  cobra.MinimumNArgs(2),
	RunE:  unsetCommand,
}

func unsetCommand(cmd *cobra.Command, args []string) error {
	name := args[0]
	keys := args[1:]
	ctx := context.Background()

	secret, err := secretsClient.Get(ctx, name)
	if err != nil {
		return err
	}

	for _, key := range keys {
		delete(secret.Data, key)
		delete(secret.Annotations, fmt.Sprintf("ksec.io/%s", key))
		fmt.Printf("Removed \"%s\" from secret \"%s\"\n", key, name)
	}

	_, err = secretsClient.Update(ctx, secret, secret.Data)
	if err != nil {
		return err
	}

	return nil
}
