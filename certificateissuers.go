package main

import (
	"github.com/risersh/controller/kubernetes/resources/secrets"
)

type NewCertificateIssuerArgs struct {
	ID     string
	Tenant string
}

// NewDeployment creates a new deployment based on requests received over RabbitMQ.
//
// Arguments:
//   - args: The arguments to create a new deployment.
//
// Returns:
//   - error: An error if the deployment could not be created.
func NewCertificateIssuer(args NewCertificateIssuerArgs) error {
	_, err := secrets.NewSecret(secrets.NewSecretArgs{
		Name:      args.ID,
		Namespace: args.Tenant,
	})
	if err != nil {
		return err
	}

	// TODO: Create kubernetes service.
	// TODO: Create kubernetes istio HTTPRoute CRD object.

	return nil
}
