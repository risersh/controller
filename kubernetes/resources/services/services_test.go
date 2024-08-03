package services

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
)

func TestNewService(t *testing.T) {
	err := NewService(NewServiceArgs{
		Name:      "test",
		Namespace: "test",
		Ports: []apiv1.ServicePort{
			{
				Name:     "test",
				Protocol: "TCP",
				Port:     80,
			},
		},
		App:  "test",
		Type: apiv1.ServiceType("ClusterIP"),
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteService(t *testing.T) {
	err := DeleteService(DeleteServiceArgs{
		Namespace: "test",
		Name:      "test",
	})

	if err != nil {
		t.Fatal(err)
	}
}
