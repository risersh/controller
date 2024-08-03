package namespaces

import (
	"context"

	"github.com/risersh/controller/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "k8s.io/api/core/v1"
)

type NewNamespaceArgs struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
}

func NewNamespace(args NewNamespaceArgs) (*apiv1.Namespace, error) {
	return kubernetes.NewNativeClient().CoreV1().Namespaces().Create(context.Background(), &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        args.Name,
			Labels:      args.Labels,
			Annotations: args.Annotations,
		},
	}, metav1.CreateOptions{})
}

func GetNamespaceByName(name string) (*apiv1.Namespace, error) {
	return kubernetes.NewNativeClient().CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
}

func DeleteNamespace(name string) error {
	return kubernetes.NewNativeClient().CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
}
