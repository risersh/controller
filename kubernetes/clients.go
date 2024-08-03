package kubernetes

import (
	"os"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewDynamicClient returns a dynamic client for accessing the Kubernetes API.
// It first tries to get the in-cluster configuration, and if that fails,
// it falls back to using the configuration from the user's kubeconfig file.
//
// Returns:
//   - dynamic.Interface: A dynamic client for accessing the Kubernetes API.
//   - error: An error if the client could not be created.
func NewDynamicClient() (dynamic.Interface, error) {
	// Load kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}

// GetClient returns a Kubernetes clientset for accessing the Kubernetes API.
// It first tries to get the in-cluster configuration, and if that fails,
// it falls back to using the configuration from the user's kubeconfig file.
//
// Returns:
//   - *kubernetes.Clientset: A Kubernetes clientset for accessing the Kubernetes API.
//   - error: An error if the client could not be created.
func NewNativeClient() (*kubernetes.Clientset, error) {
	// Temporary until we get auth tokens working.
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// If the environment is not kubernetes, use the kubeconfig file.
	config, err := rest.InClusterConfig()
	if os.Getenv("ENVIRONMENT") != "kubernetes" {
		if err != nil {
			config, err = clientcmd.BuildConfigFromFlags("", home+"/.kube/config")
			if err != nil {
				panic(err.Error())
			}
		}
	}

	// Create a Kubernetes clientset for accessing the Kubernetes API.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
