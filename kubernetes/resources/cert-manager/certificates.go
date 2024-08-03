package certmanager

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/risersh/controller/kubernetes"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type IssuerEnvironment string

const (
	IssuerEnvironmentProd    IssuerEnvironment = "prod"
	IssuerEnvironmentStaging IssuerEnvironment = "staging"
)

var CertificatesResource = schema.GroupVersionResource{
	Group:    "cert-manager.io/v1",
	Resource: "certificates",
}

func ListCertificates() (*unstructured.UnstructuredList, error) {
	client, err := kubernetes.NewDynamicClient()
	if err != nil {
		return nil, err
	}

	certificates, err := client.Resource(CertificatesResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

func NewCertificate(name string, namespace string, hostnames []string, env IssuerEnvironment) (*unstructured.Unstructured, error) {
	client, err := kubernetes.NewDynamicClient()
	if err != nil {
		return nil, err
	}

	un := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1",
			"kind":       "Certificate",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"dnsNames": hostnames,
				"issuerRef": map[string]interface{}{
					"kind": "ClusterIssuer",
					"name": fmt.Sprintf("letsencrypt-%s", env),
				},
				"secretName": name,
			},
		},
	}

	res, err := client.Resource(CertificatesResource).Namespace(namespace).Create(context.Background(), un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err = wait.PollUntilContextTimeout(ctx, 5*time.Second, 2*time.Minute, true, func(ctx context.Context) (bool, error) {
		cert, err := client.Resource(CertificatesResource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		status, found, err := unstructured.NestedString(cert.Object, "status", "conditions", "0", "type")
		return found && status == "Ready", err
	})

	if err != nil {
		return nil, fmt.Errorf("certificate not issued within timeout: %v", err)
	}

	return res, nil
}
