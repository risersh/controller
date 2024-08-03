package kubernetes

import (
	"fmt"
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/risersh/controller/conf"
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
func NewNativeClient() *kubernetes.Clientset {
	var config *rest.Config
	var err error

	// If the environment is containerized, use the in-cluster configuration
	// otherwise, use the kubeconfig file from the user's home directory.
	if conf.Config.Environment.Containerized {
		config, err = rest.InClusterConfig()
		if err != nil {
			multilog.Fatal("kubernetes.clients", "get in-cluster configuration", map[string]interface{}{
				"error": err,
			})
		}
	} else {
		// Temporary until we get auth tokens working.
		home, err := os.UserHomeDir()
		if err != nil {
			multilog.Fatal("kubernetes.clients", "get user home directory", map[string]interface{}{
				"error": err,
			})
		}

		config, err = clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home))
		if err != nil {
			multilog.Fatal("kubernetes.clients", "new native client", map[string]interface{}{
				"error": err,
			})
		}
	}

	// Create a Kubernetes clientset for accessing the Kubernetes API using the
	// configuration we just created above.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		multilog.Fatal("kubernetes.clients", "new clientset", map[string]interface{}{
			"error": err,
		})
	}

	return clientset
}
