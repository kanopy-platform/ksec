package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"gopkg.in/urfave/cli.v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type keyAnnotation struct {
	UpdatedBy   string `json:"updatedBy"`
	LastUpdated string `json:"lastUpdated"`
}

type Config struct {
	KubeConfig      clientcmd.ClientConfig
	SecretInterface apiv1.SecretInterface
	User            string
}

var cfg Config

func listSecrets(ctx *cli.Context) error {
	secrets, err := cfg.SecretInterface.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Println("NAME")
	for _, secret := range secrets.Items {
		fmt.Println(secret.Name)
	}
	return nil
}

func createSecret(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	name := ctx.Args().Get(0)

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if _, err := cfg.SecretInterface.Create(&secret); err != nil {
		return err
	}

	fmt.Printf("Created secret \"%s\"\n", name)
	return nil
}

func deleteSecrets(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("No arguments specified")
	}

	for _, secret := range ctx.Args() {
		if err := cfg.SecretInterface.Delete(secret, &metav1.DeleteOptions{}); err != nil {
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

	secret, err := cfg.SecretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
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

		annotation := keyAnnotation{}
		if err := json.Unmarshal(jsonAnnotation, &annotation); err != nil {
			return err
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s", key, value, annotation.UpdatedBy, annotation.LastUpdated))
	}
	output_tabular(lines)

	return nil
}

func setSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := cfg.SecretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
	if err != nil {
		return err
	}

	annotation := keyAnnotation{
		UpdatedBy:   cfg.User,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	if secret.ObjectMeta.Annotations == nil {
		secret.ObjectMeta.Annotations = make(map[string]string)
	}

	data := ctx.Args().Get(1)

	for _, item := range strings.Split(data, ",") {
		split := strings.SplitN(item, "=", 2)
		if len(split) != 2 {
			return fmt.Errorf("Data is not formatted correctly: %s", item)
		}
		secret.Data[split[0]] = []byte(split[1])
		jsonAnnotations, err := json.Marshal(annotation)
		if err != nil {
			return err
		}
		secret.ObjectMeta.Annotations[fmt.Sprintf("ksec.io/%s", split[0])] = string(jsonAnnotations)
	}

	_, err = cfg.SecretInterface.Update(secret)
	if err != nil {
		return err
	}

	return nil
}

func unsetSecretKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := cfg.SecretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
	if err != nil {
		return err
	}

	keys := ctx.Args().Get(1)

	for _, key := range strings.Split(keys, ",") {
		delete(secret.Data, key)
		fmt.Printf("Removed \"%s\" from secret \"%s\"\n", key, secret.Name)
	}

	_, err = cfg.SecretInterface.Update(secret)
	if err != nil {
		return err
	}

	return nil
}

func pushKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := cfg.SecretInterface.Get(ctx.Args().Get(1), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	file, err := os.Open(ctx.Args().Get(0))
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

		secret.Data[split[0]] = []byte(split[1])
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	_, err = cfg.SecretInterface.Update(secret)
	if err != nil {
		return err
	}

	return nil
}

func pullKeys(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return fmt.Errorf("Incorrect number of arguments")
	}

	secret, err := cfg.SecretInterface.Get(ctx.Args().Get(0), metav1.GetOptions{})
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

func output_tabular(lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}

func main() {
	app := cli.NewApp()
	app.Usage = "A tool for managing Kubernetes Secret data"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "namespace, n",
			Usage: "Operate in a specific `NAMESPACE` (Defaults to the current kubeconfig namespace)",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all secrets in a namespace",
			Action: func(ctx *cli.Context) error {
				return listSecrets(ctx)
			},
		},
		{
			Name:  "create",
			Usage: "Create a Secret",
			Action: func(ctx *cli.Context) error {
				return createSecret(ctx)
			},
		},
		{
			Name:  "delete",
			Usage: "Delete a Secret",
			Action: func(ctx *cli.Context) error {
				return deleteSecrets(ctx)
			},
		},
		{
			Name:  "get",
			Usage: "Get values from a Secret",
			Action: func(ctx *cli.Context) error {
				return getSecretKeys(ctx)
			},
		},
		{
			Name:  "set",
			Usage: "Set values in a Secret",
			Action: func(ctx *cli.Context) error {
				return setSecretKeys(ctx)
			},
		},
		{
			Name:  "unset",
			Usage: "Unset values in a Secret",
			Action: func(ctx *cli.Context) error {
				return unsetSecretKeys(ctx)
			},
		},
		{
			Name:  "push",
			Usage: "Push values from a .env file into a Secret",
			Action: func(ctx *cli.Context) error {
				return pushKeys(ctx)
			},
		},
		{
			Name:  "pull",
			Usage: "Pull values from a Secret into a .env file",
			Action: func(ctx *cli.Context) error {
				return pullKeys(ctx)
			},
		},
	}

	app.Before = func(ctx *cli.Context) error {
		// set KubeConfig global
		cfg.KubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)

		// set SecretInterface global
		var namespace string
		config, err := cfg.KubeConfig.ClientConfig()
		if err != nil {
			return err
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}
		if ctx.String("namespace") != "" {
			namespace = ctx.String("namespace")
		} else {
			namespace, _, err = cfg.KubeConfig.Namespace()
			if err != nil {
				return err
			}
		}
		cfg.SecretInterface = client.CoreV1().Secrets(namespace)

		// set User global
		rawConfig, err := cfg.KubeConfig.RawConfig()
		if err != nil {
			return err
		}
		cfg.User = rawConfig.Contexts[rawConfig.CurrentContext].AuthInfo
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
