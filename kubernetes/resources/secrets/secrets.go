package secrets

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/risersh/controller/kubernetes"
)

type NewSecretArgs struct {
	Name      string
	Namespace string
	Data      map[string][]byte
	Labels    map[string]string
}

func NewSecret(args NewSecretArgs) (*corev1.Secret, error) {
	client, err := kubernetes.NewNativeClient()
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      args.Name,
			Namespace: args.Namespace,
			Labels:    args.Labels,
		},
		Data: args.Data,
		Type: corev1.SecretTypeOpaque,
	}

	res, err := client.CoreV1().Secrets(args.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetSecret(name, namespace string) (*corev1.Secret, error) {
	client, err := kubernetes.NewNativeClient()
	if err != nil {
		return nil, err
	}

	return client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func GetSecretsByLabels(namespace string, matchLabels map[string]string) ([]corev1.Secret, error) {
	client, err := kubernetes.NewNativeClient()
	if err != nil {
		return nil, err
	}

	labelSelector := labels.SelectorFromSet(matchLabels).String()
	res, err := client.CoreV1().Secrets(namespace).List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: labelSelector,
		},
	)
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func DeleteSecret(name, namespace string) error {
	client, err := kubernetes.NewNativeClient()
	if err != nil {
		return err
	}

	return client.CoreV1().Secrets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
