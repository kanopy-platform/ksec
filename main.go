package main

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var secretInterface apiv1.SecretInterface

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

func list() {
	secrets, err := secretInterface.List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, secret := range secrets.Items {
		fmt.Printf("%s\t%s\n", secret.Name, secret.Namespace)
	}
}

func main() {
	list()
}
