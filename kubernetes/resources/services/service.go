package services

import (
	"context"
	"log"

	"github.com/risersh/controller/kubernetes"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NewServiceArgs struct {
	Name      string
	App       string
	Namespace string
	Labels    map[string]string
	Ports     []apiv1.ServicePort
	Type      apiv1.ServiceType
}

func NewService(args NewServiceArgs) error {

	client, _ := kubernetes.NewNativeClient()
	servicesClient := client.CoreV1().Services(args.Namespace)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      args.Name,
			Namespace: args.Namespace,
			Labels:    args.Labels,
		},
		Spec: apiv1.ServiceSpec{
			Type: args.Type,
			Selector: map[string]string{
				"app": args.App,
			},
			Ports: args.Ports,
		},
	}
	_, err := servicesClient.Create(context.Background(), service, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	log.Printf("Created service %q in namespace %q \n", service.GetObjectMeta().GetName(), service.GetObjectMeta().GetNamespace())

	return nil

}
