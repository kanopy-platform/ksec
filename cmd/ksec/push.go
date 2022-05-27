package main

import (
	"bufio"
	"context"
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
	ctx := context.Background()

	file, err := os.Open(fileArg)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := readSecretData(file)
	if err != nil {
		return err
	}

	_, err = secretsClient.Upsert(ctx, secretName, data)
	if err != nil {
		return err
	}
	return nil
}

func readSecretData(reader io.Reader) (map[string][]byte, error) {
	data := map[string][]byte{}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		split := strings.SplitN(scanner.Text(), "=", 2)

		if len(split) > 1 {
			data[split[0]] = []byte(split[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
