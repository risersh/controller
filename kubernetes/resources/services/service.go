package services

import (
	"context"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/risersh/controller/kubernetes"
	"github.com/risersh/util/variables"
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
	res, err := kubernetes.NewNativeClient().CoreV1().Services(args.Namespace).Create(context.Background(), &apiv1.Service{
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
	}, metav1.CreateOptions{})
	if err != nil {
		multilog.Error("kubernetes.services.create", "error creating service", map[string]interface{}{
			"error": err,
		})
		return err
	}

	multilog.Info("kubernetes.services.create", "created service", map[string]interface{}{
		"namespace": res.GetObjectMeta().GetNamespace(),
		"name":      res.GetObjectMeta().GetName(),
	})

	return nil
}

type DeleteServiceArgs struct {
	Namespace string
	Name      string
}

func DeleteService(args DeleteServiceArgs) error {
	if err := kubernetes.NewNativeClient().CoreV1().Services(args.Namespace).Delete(context.Background(), args.Name, metav1.DeleteOptions{
		PropagationPolicy: variables.ToPtr(metav1.DeletePropagationForeground),
	}); err != nil {
		multilog.Error("kubernetes.services.delete", "error deleting service", map[string]interface{}{
			"error": err,
		})
		return err
	}

	multilog.Info("kubernetes.services.delete", "delete service", map[string]interface{}{
		"namespace": args.Namespace,
		"name":      args.Name,
	})

	return nil
}
