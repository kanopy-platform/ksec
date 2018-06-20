package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset [secret] [key...]",
	Short: "Unset values in a Secret",
	Args:  cobra.MinimumNArgs(2),
	Run:   unsetCommand,
}

func init() {
	RootCmd.AddCommand(unsetCmd)
}

func unsetCommand(cmd *cobra.Command, args []string) {
	name := args[0]
	keys := args[1:]

	secret, err := secretsClient.Get(name)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, key := range keys {
		delete(secret.Data, key)
		delete(secret.Annotations, fmt.Sprintf("ksec.io/%s", key))
		fmt.Printf("Removed \"%s\" from secret \"%s\"\n", key, name)
	}

	_, err = secretsClient.Update(secret, secret.Data)
	if err != nil {
		log.Fatal(err.Error())
	}
}
