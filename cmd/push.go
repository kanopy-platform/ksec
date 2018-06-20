package cmd

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [file] [secret]",
	Short: "Push values from a .env file into a Secret",
	Args:  cobra.ExactArgs(2),
	Run:   pushCommand,
}

func init() {
	RootCmd.AddCommand(pushCmd)
}

func pushCommand(cmd *cobra.Command, args []string) {
	fileArg := args[0]
	name := args[1]
	data := make(map[string][]byte)

	file, err := os.Open(fileArg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		split := strings.SplitN(text, "=", 2)
		data[split[0]] = []byte(split[1])
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err.Error())
	}

	_, err = secretsClient.Upsert(name, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}
