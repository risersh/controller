package certmanager

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/risersh/controller/kubernetes"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var CertificatesResource = schema.GroupVersionResource{
	Group:    "cert-manager.io",
	Version:  "v1",
	Resource: "certificates",
}

func ListCertificates() (*unstructured.UnstructuredList, error) {
	client, err := kubernetes.NewClient()
	if err != nil {
		return nil, err
	}

	certificates, err := client.Resource(CertificatesResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
