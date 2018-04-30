package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/colinhoglund/helm-k8s-secrets/pkg/models"
	"gopkg.in/urfave/cli.v1"
)

func listCommand(ctx *cli.Context) error {
	secrets, err := secretsClient.List()
	if err != nil {
		return err
	}

	fmt.Println("NAME")
	for _, secret := range secrets.Items {
		fmt.Println(secret.Name)
	}
	return nil
}

func createCommand(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	name := ctx.Args().Get(0)
	_, err := secretsClient.Create(name)
	if err != nil {
		return err
	}

	fmt.Printf("Created secret \"%s\"\n", name)
	return nil
}

func deleteSecrets(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("No arguments specified")
	}

	for _, name := range ctx.Args() {
		if err := secretsClient.Delete(name); err != nil {
			return err
		}

		fmt.Printf("Deleted secret \"%s\"\n", name)
	}
	return nil
}

func getSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := secretsClient.Get(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	lines := []string{"KEY\tVALUE\tUSER\tUPDATED"}
	for key, value := range secret.Data {
		rawAnnotation := secret.Annotations[fmt.Sprintf("ksec.io/%s", key)]

		var jsonAnnotation []byte
		if rawAnnotation == "" {
			jsonAnnotation = []byte(`{"updatedBy": "", "lastUpdated": ""}`)
		} else {
			jsonAnnotation = []byte(rawAnnotation)
		}

		annotation := models.KeyAnnotation{}
		if err := json.Unmarshal(jsonAnnotation, &annotation); err != nil {
			return err
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s", key, value, annotation.UpdatedBy, annotation.LastUpdated))
	}
	outputTabular(lines)

	return nil
}

func setSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	// annotation := models.KeyAnnotation{
	// 	UpdatedBy:   secretsClient.AuthInfo,
	// 	LastUpdated: time.Now().Format(time.RFC3339),
	// }

	// if secret.ObjectMeta.Annotations == nil {
	// 	secret.ObjectMeta.Annotations = make(map[string]string)
	// }

	name := ctx.Args().Get(0)
	dataArgs := ctx.Args().Get(1)
	data := make(map[string][]byte)

	for _, item := range strings.Split(dataArgs, ",") {
		split := strings.SplitN(item, "=", 2)
		if len(split) != 2 {
			return fmt.Errorf("Data is not formatted correctly: %s", item)
		}
		data[split[0]] = []byte(split[1])
		// jsonAnnotations, err := json.Marshal(annotation)
		// if err != nil {
		// 	return err
		// }
		// secret.ObjectMeta.Annotations[fmt.Sprintf("ksec.io/%s", split[0])] = string(jsonAnnotations)
	}

	_, err := secretsClient.Update(name, data)
	if err != nil {
		return err
	}

	return nil
}

func unsetSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	name := ctx.Args().Get(0)
	keys := ctx.Args().Get(1)

	secret, err := secretsClient.Get(name)
	if err != nil {
		return err
	}

	for _, key := range strings.Split(keys, ",") {
		delete(secret.Data, key)
		fmt.Printf("Removed \"%s\" from secret \"%s\"\n", key, name)
	}

	_, err = secretsClient.Update(name, secret.Data)
	if err != nil {
		return err
	}

	return nil
}

func pushKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	fileArg := ctx.Args().Get(0)
	name := ctx.Args().Get(1)
	data := make(map[string][]byte)

	file, err := os.Open(fileArg)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		split := strings.Split(text, "=")

		if len(split) != 2 {
			return fmt.Errorf("Incorrectly formatted environment variable: %s", text)
		}

		data[split[0]] = []byte(split[1])
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	_, err = secretsClient.Update(name, data)
	if err != nil {
		return err
	}

	return nil
}

func pullKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	name := ctx.Args().Get(0)

	secret, err := secretsClient.Get(name)
	if err != nil {
		return err
	}

	file, err := os.Create(ctx.Args().Get(1))
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

func outputTabular(lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
