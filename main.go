package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"gopkg.in/urfave/cli.v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	secretInterface apiv1.SecretInterface
	namespace       string
)

func init() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	var err error
	namespace, _, err = kubeConfig.Namespace()
	if err != nil {
		log.Fatal(err.Error())
	}

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	secretInterface = client.CoreV1().Secrets(namespace)
}

func listSecrets() error {
	secrets, err := secretInterface.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	lines := []string{"NAME\tNAMESPACE"}
	for _, secret := range secrets.Items {
		lines = append(lines, fmt.Sprintf("%s\t%s", secret.Name, namespace))
	}
	output_tabular(lines)
	return nil
}

func createSecret(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	name := ctx.Args().Get(0)

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if _, err := secretInterface.Create(&secret); err != nil {
		return err
	}

	fmt.Printf("Created secret \"%s\"\n", name)
	return nil
}

func deleteSecret(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("No arguments specified")
	}

	for _, secret := range ctx.Args() {
		if err := secretInterface.Delete(secret, &metav1.DeleteOptions{}); err != nil {
			return err
		}

		fmt.Printf("Deleted secret \"%s\"\n", secret)
	}
	return nil
}

func getSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := secretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
	if err != nil {
		return err
	}

	lines := []string{"KEY\tVALUE"}
	for key, value := range secret.Data {
		lines = append(lines, fmt.Sprintf("%s\t%s", key, value))
	}
	output_tabular(lines)

	return nil
}

func setSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := secretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	data := ctx.Args().Get(1)

	for _, item := range strings.Split(data, ",") {
		split := strings.SplitN(item, "=", 2)
		if len(split) != 2 {
			return fmt.Errorf("Data is not formatted correctly: %s", item)
		}
		secret.Data[split[0]] = []byte(split[1])
	}

	_, err = secretInterface.Update(secret)
	if err != nil {
		return err
	}

	return nil
}

func unsetSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := secretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
	if err != nil {
		return err
	}

	keys := ctx.Args().Get(1)

	for _, key := range strings.Split(keys, ",") {
		delete(secret.Data, key)
		fmt.Printf("Removed \"%s\" from secret \"%s\"\n", key, secret.Name)
	}

	_, err = secretInterface.Update(secret)
	if err != nil {
		return err
	}

	return nil
}

func output_tabular(lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}

func loadKeys(ctx *cli.Context) error {
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = "A tool for managing Kubernetes Secret data"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all secrets in a namespace",
			Action: func(ctx *cli.Context) error {
				return listSecrets()
			},
		},
		{
			Name:  "create",
			Usage: "Create a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				return createSecret(ctx)
			},
		},
		{
			Name:  "delete",
			Usage: "Delete a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				return deleteSecret(ctx)
			},
		},
		{
			Name:  "get",
			Usage: "Get values from a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				return getSecretKeys(ctx)
			},
		},
		{
			Name:  "set",
			Usage: "Set values in a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				return setSecretKeys(ctx)
			},
		},
		{
			Name:  "unset",
			Usage: "Unset values in a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				return unsetSecretKeys(ctx)
			},
		},
		{
			Name:  "load",
			Usage: "Load values from a env file into a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				return loadKeys(ctx)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
