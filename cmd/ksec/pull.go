package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [secret] [file]",
	Short: "Pull values from a Secret into a .env file",
	Args:  cobra.ExactArgs(2),
	RunE:  pullCommand,
}

func pullCommand(cmd *cobra.Command, args []string) error {
	name := args[0]
	ctx := context.Background()

	secret, err := secretsClient.Get(ctx, name)
	if err != nil {
		return err
	}

	file, err := os.Create(args[1])
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range secret.Data {
		_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return err
		}
	}

	file.Sync()
	return nil
}
