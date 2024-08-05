package main

import (
	"github.com/risersh/controller/kubernetes/resources/deployments"
	appsv1 "k8s.io/api/apps/v1"
)

type NewDeploymentArgs struct {
	ID       string
	Tenant   string
	Replicas int
	Image    string
}

// NewDeployment creates a new deployment based on requests received over RabbitMQ.
//
// Arguments:
//   - args: The arguments to create a new deployment.
//
// Returns:
//   - error: An error if the deployment could not be created.
func NewDeployment(args NewDeploymentArgs) (*appsv1.Deployment, []error) {
	deployment, err := deployments.NewDeployment(deployments.NewDeploymentArgs{
		Name:      args.ID,
		Namespace: args.Tenant,
		Replicas:  1,
		Image:     args.Image,
	})
	if len(err) > 0 {
		return nil, err
	}

	// TODO: Create kubernetes service.
	// TODO: Create kubernetes istio HTTPRoute CRD object.

	return deployment, []error{}
}
