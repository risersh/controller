package deployments

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
)

func TestNewDeployment(t *testing.T) {
	_, err := NewDeployment(NewDeploymentArgs{
		Name:      "test",
		Namespace: "test",
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
	})
	if err != nil {
		t.Fatal(err)
	}
}
