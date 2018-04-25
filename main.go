package main

import (
	"fmt"
	"log"
	"os"
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

func listSecrets() {
	secrets, err := secretInterface.List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	lines := []string{"Name\tNamespace", "----\t---------"}
	for _, secret := range secrets.Items {
		lines = append(lines, fmt.Sprintf("%s\t%s", secret.Name, namespace))
	}
	output_tabular(lines)
}

func createSecret(ctx *cli.Context) {
	if len(ctx.Args()) != 1 {
		log.Fatal("ERROR: Incorrect number of arguments")
	}

	name := ctx.Args().Get(0)

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := secretInterface.Create(&secret)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Created %s\n", name)
}

func deleteSecret(ctx *cli.Context) {
	if len(ctx.Args()) < 1 {
		log.Fatal("No arguments specified")
	}

	for _, secret := range ctx.Args() {
		if err := secretInterface.Delete(secret, &metav1.DeleteOptions{}); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Deleted %s\n", secret)
	}
}

func getSecretKeys(ctx *cli.Context) {
	return
}

func setSecretKeys(ctx *cli.Context) {
	return
}

func output_tabular(lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}

func main() {
	app := cli.NewApp()
	app.Usage = "A tool managing Kubernetes Secret data"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all secrets in a namespace",
			Action: func(ctx *cli.Context) error {
				listSecrets()
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Create a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				createSecret(ctx)
				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "Delete a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				deleteSecret(ctx)
				return nil
			},
		},
		{
			Name:  "get",
			Usage: "Get values from a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				getSecretKeys(ctx)
				return nil
			},
		},
		{
			Name:  "set",
			Usage: "Set values in a Kubernetes Secret",
			Action: func(ctx *cli.Context) error {
				setSecretKeys(ctx)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
