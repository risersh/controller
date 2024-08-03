package certmanager

import (
	"context"

	"github.com/risersh/controller/conf"
	"github.com/risersh/controller/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type IssuerSolverType string

const (
	IssuerSolverTypeDNS  IssuerSolverType = "dns"
	IssuerSolverTypeHTTP IssuerSolverType = "http"
)

var IssuerResource = schema.GroupVersionResource{
	Group:    "cert-manager.io",
	Resource: "issuers",
	Version:  "v1",
}

type NewIssuerArgs struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Solver    IssuerSolverType
}

func NewIssuer(args NewIssuerArgs) (*unstructured.Unstructured, error) {
	client, err := kubernetes.NewDynamicClient()
	if err != nil {
		return nil, err
	}

	var solvers []map[string]interface{}
	switch args.Solver {
	case IssuerSolverTypeDNS:
		solver := map[string]interface{}{
			"email": conf.Config.Certificates.Email,
			"apiKeySecretRef": map[string]interface{}{
				"name": "cloudflare-api-key",
				"key":  "api-key",
			},
		}
		solvers = append(solvers, map[string]interface{}{
			"dns01": solver,
		})
	case IssuerSolverTypeHTTP:
		solver := map[string]interface{}{
			"ingress": map[string]interface{}{
				"ingressClassName": "istio",
			},
		}
		solvers = append(solvers, map[string]interface{}{
			"http01": solver,
		})
	}

	un := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1",
			"kind":       "Issuer",
			"metadata": map[string]interface{}{
				"name":      args.Name,
				"namespace": args.Namespace,
				"labels":    args.Labels,
			},
			"spec": map[string]interface{}{
				"acme": map[string]interface{}{
					"server": conf.Config.Certificates.Server,
					"email":  conf.Config.Certificates.Email,
					"privateKeySecretRef": map[string]interface{}{
						"name": args.Name,
					},
					"solvers": solvers,
				},
			},
		},
	}

	res, err := client.Resource(IssuerResource).Namespace(args.Namespace).Create(context.Background(), un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
