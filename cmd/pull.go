package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [secret] [file]",
	Short: "Pull values from a Secret into a .env file",
	Args:  cobra.ExactArgs(2),
	Run:   pullCommand,
}

func init() {
	RootCmd.AddCommand(pullCmd)
}

func pullCommand(cmd *cobra.Command, args []string) {
	name := args[0]

	secret, err := secretsClient.Get(name)
	if err != nil {
		log.Fatal(err.Error())
	}

	file, err := os.Create(args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	for key, value := range secret.Data {
		_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	file.Sync()
}
