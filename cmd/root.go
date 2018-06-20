package cmd

import (
	"fmt"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/colinhoglund/ksec/pkg/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var secretsClient *models.SecretsClient

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "ksec",
	Short:   "A tool for managing Kubernetes Secret data",
	Version: "0.1.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		secretsClient, err = models.NewSecretsClient(viper.GetString("namespace"))
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global options
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (Default: $HOME/.ksec.yaml)")
	RootCmd.PersistentFlags().StringP("namespace", "n", "", "Operate in a specific NAMESPACE (Default: current kubeconfig namespace)")

	// setup viper config
	viper.BindPFlags(RootCmd.PersistentFlags())
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
