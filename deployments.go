package main

import "github.com/risersh/controller/kubernetes/resources/deployments"

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
func NewDeployment(args NewDeploymentArgs) error {
	_, err := deployments.NewDeployment(deployments.NewDeploymentArgs{
		Name:      args.ID,
		Namespace: args.Tenant,
		Replicas:  1,
		Image:     args.Image,
	})
	if err != nil {
		return err
	}

	// TODO: Create kubernetes service.
	// TODO: Create kubernetes istio HTTPRoute CRD object.

	return nil
}
