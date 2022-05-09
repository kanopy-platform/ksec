package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [file] [secret]",
	Short: "Push values from a .env file into a Secret",
	Args:  cobra.ExactArgs(2),
	RunE:  pushCommand,
}

func pushCommand(cmd *cobra.Command, args []string) error {
	fileArg := args[0]
	secretName := args[1]
	data := make(map[string][]byte)

	file, err := os.Open(fileArg)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := scanSecretData(file, data); err != nil {
		return err
	}

	_, err = secretsClient.Upsert(secretName, data)
	if err != nil {
		return err
	}
	return nil
}

func scanSecretData(reader io.Reader, data map[string][]byte) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		text := scanner.Text()
		split := strings.SplitN(text, "=", 2)
		data[split[0]] = []byte(split[1])
	}

	return scanner.Err()
}
