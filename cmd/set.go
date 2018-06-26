package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [secret] [key=value...]",
	Short: "Set values in a Secret",
	Args:  cobra.MinimumNArgs(2),
	Run:   setCommand,
}

func setCommand(cmd *cobra.Command, args []string) {
	name := args[0]
	dataArgs := args[1:]
	data := make(map[string][]byte)

	for _, item := range dataArgs {
		split := strings.SplitN(item, "=", 2)
		if len(split) != 2 {
			log.Fatal("Data is not formatted correctly: ", item)
		}
		data[split[0]] = []byte(split[1])
	}

	_, err := secretsClient.Upsert(name, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}
