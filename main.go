package main

import (
	"fmt"
	"log"
	"os"

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

	namespace, _, err := kubeConfig.Namespace()
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

	for _, secret := range secrets.Items {
		fmt.Printf("%s\n", secret.Name)
	}
}

func createSecret(c *cli.Context) {
	if len(c.Args()) != 1 {
		log.Fatal("ERROR: Incorrect number of arguments")
	}

	name := c.Args().Get(0)

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

func deleteSecret(c *cli.Context) {
	if len(c.Args()) < 1 {
		log.Fatal("No arguments specified")
	}

	for _, secret := range c.Args() {
		if err := secretInterface.Delete(secret, &metav1.DeleteOptions{}); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Deleted %s\n", secret)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all secrets in a namespace",
			Action: func(c *cli.Context) error {
				listSecrets()
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Create a Kubernetes Secret",
			Action: func(c *cli.Context) error {
				createSecret(c)
				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "Delete a Kubernetes Secret",
			Action: func(c *cli.Context) error {
				deleteSecret(c)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
