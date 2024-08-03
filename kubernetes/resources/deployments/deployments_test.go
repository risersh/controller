package deployments

import (
	"testing"
	"time"

	"github.com/risersh/controller/conf"
	"github.com/risersh/controller/kubernetes/resources/namespaces"
	"github.com/risersh/controller/monitoring"
	apiv1 "k8s.io/api/core/v1"
)

func TestNewDeployment(t *testing.T) {
	conf.Init()
	monitoring.Setup([]monitoring.LoggerType{
		monitoring.LoggerTypeConsole,
	})

	namespace, err := namespaces.NewNamespace(namespaces.NewNamespaceArgs{
		Name: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer namespaces.DeleteNamespace(namespace.Name)

	_, errs := NewDeployment(NewDeploymentArgs{
		Name:      "test",
		Namespace: namespace.Name,
		Replicas:  1,
		Image:     "nginx",
		Labels: map[string]string{
			"app": "test",
		},
		Ports: []apiv1.ContainerPort{
			{
				ContainerPort: 80,
			},
		},
		Resources: ResourcesArgs{
			Requests: ResourceArgs{
				CPU:    "100m",
				Memory: "100Mi",
			},
			Limits: ResourceArgs{
				CPU:    "100m",
				Memory: "100Mi",
			},
		},
		EnvVars: []apiv1.EnvVar{
			{
				Name:  "FOO",
				Value: "bar",
			},
		},
		Timeout: 15 * time.Second,
	})
	if len(errs) > 0 {
		t.Fatal(errs)
	}
}

// func TestDeleteDeployment(t *testing.T) {
// 	err := DeleteDeployment(DeleteDeploymentArgs{
// 		Name:      "test",
// 		Namespace: "test",
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
