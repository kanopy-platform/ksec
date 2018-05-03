package main

import (
	"log"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/colinhoglund/helm-k8s-secrets/pkg/models"
	"gopkg.in/urfave/cli.v1"
)

var secretsClient *models.SecretsClient

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
				return listCommand(ctx)
			},
		},
		{
			Name:  "create",
			Usage: "Create a Secret",
			Action: func(ctx *cli.Context) error {
				return createCommand(ctx)
			},
		},
		{
			Name:  "delete",
			Usage: "Delete a Secret",
			Action: func(ctx *cli.Context) error {
				return deleteCommand(ctx)
			},
		},
		{
			Name:  "get",
			Usage: "Get values from a Secret",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Show extra metadata",
				},
			},
			Action: func(ctx *cli.Context) error {
				return getCommand(ctx)
			},
		},
		{
			Name:  "set",
			Usage: "Set values in a Secret",
			Action: func(ctx *cli.Context) error {
				return setCommand(ctx)
			},
		},
		{
			Name:  "unset",
			Usage: "Unset values in a Secret",
			Action: func(ctx *cli.Context) error {
				return unsetCommand(ctx)
			},
		},
		{
			Name:  "push",
			Usage: "Push values from a .env file into a Secret",
			Action: func(ctx *cli.Context) error {
				return pushCommand(ctx)
			},
		},
		{
			Name:  "pull",
			Usage: "Pull values from a Secret into a .env file",
			Action: func(ctx *cli.Context) error {
				return pullCommand(ctx)
			},
		},
	}

	app.Before = func(ctx *cli.Context) error {
		var err error
		secretsClient, err = models.NewSecretsClient(ctx.String("namespace"))
		if err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
