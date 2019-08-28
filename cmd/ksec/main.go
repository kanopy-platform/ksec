package main

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/10gen-ops/ksec/pkg/models"
	"github.com/10gen-ops/ksec/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var secretsClient *models.SecretsClient

func initRootCmd(rootCmd *cobra.Command) {
	cobra.OnInitialize(initConfig)

	// global options
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (Default: $HOME/.ksec.yaml)")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Operate in a specific NAMESPACE (Default: current kubeconfig namespace)")

	// setup viper config
	viper.BindPFlags(rootCmd.PersistentFlags())

	// subcommands without extra options
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)

	// subcommands with extra options
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("verbose", "v", false, "Show extra metadata")

	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("yes", "y", false, "do not ask for confirmation")

	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Show all secrets (Default: Opaque only)")

	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(bashCompletionCmd)
	completionCmd.AddCommand(zshCompletionCmd)
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ksec",
		Short:   "A tool for managing Kubernetes Secret data",
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			secretsClient, err = models.NewSecretsClient(viper.GetString("namespace"))
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	initRootCmd(rootCmd)
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Search config in home directory with name ".ksec" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ksec")
	}

	viper.SetEnvPrefix("KSEC")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
