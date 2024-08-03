package deployments

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/risersh/controller/kubernetes"
	"github.com/risersh/util/variables"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NewDeploymentArgs struct {
	Name         string                `json:"name"`
	Namespace    string                `json:"namespace"`
	Replicas     int32                 `json:"replicas"`
	Image        string                `json:"image"`
	Labels       map[string]string     `json:"labels"`
	NodeSelector string                `json:"nodeSelector"`
	Ports        []apiv1.ContainerPort `json:"ports"`
	EnvVars      []apiv1.EnvVar        `json:"envVars"`
	Resources    ResourcesArgs         `json:"resources"`
	Timeout      time.Duration         `json:"timeout"`
}

type ResourcesArgs struct {
	Requests ResourceArgs `json:"requests"`
	Limits   ResourceArgs `json:"limits"`
}

type ResourceArgs struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type DeleteDeploymentArgs struct {
	Namespace string
	Name      string
}

func NewDeployment(args NewDeploymentArgs) (*appsv1.Deployment, []error) {
	multilog.Info("kubernetes.deployments.new", "Creating kubernetes deployment", map[string]interface{}{
		"args": args,
	})

	// Create a new deployment.
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      args.Name,
			Namespace: args.Namespace,
			Labels:    args.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &args.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: args.Labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: args.Labels,
				},
				Spec: apiv1.PodSpec{
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: "ghcr",
						},
					},
					TerminationGracePeriodSeconds: &[]int64{10}[0],
					Containers: []apiv1.Container{
						{
							Name:  "default",
							Image: args.Image,
							Ports: args.Ports,
							Env:   args.EnvVars,
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									"cpu":    resource.MustParse(args.Resources.Requests.CPU),
									"memory": resource.MustParse(args.Resources.Requests.Memory),
								},
								Limits: apiv1.ResourceList{
									"cpu":    resource.MustParse(args.Resources.Limits.CPU),
									"memory": resource.MustParse(args.Resources.Limits.Memory),
								},
							},
						},
					},
				},
			},
		},
	}

	// If a node selector is provided, set it on the deployment.
	if args.NodeSelector != "" {
		deployment.Spec.Template.Spec.NodeSelector = map[string]string{
			"role": args.NodeSelector,
		}
	}

	client := kubernetes.NewNativeClient()

	// Create the deployment in the cluster.
	result, err := client.AppsV1().Deployments(args.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		multilog.Error("kubernetes.deployments.new", "error creating kubernetes deployment", map[string]interface{}{
			"error": err,
		})
		return nil, []error{err}
	}

	multilog.Info("kubernetes.deployments.new", "created kubernetes deployment", map[string]interface{}{
		"name": result.GetObjectMeta().GetName(),
	})

	// Create a context with a timeout for the api call to get the deployment.
	// This is to prevent the api call from blocking indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), args.Timeout)
	defer cancel()

	// Wait for the deployment to be ready by polling the deployment status.
	// This is a blocking call that will return an error if the deployment
	// is not ready within the timeout.
	errs := kubernetes.WaitForResourceConditions(kubernetes.WaitForResourceConditionArgs{
		Timeout: args.Timeout,
		Evaluator: func() (bool, error) {
			deployment, err := client.AppsV1().Deployments(args.Namespace).Get(ctx, args.Name, metav1.GetOptions{})
			if err != nil {
				multilog.Error("kubernetes.deployments.new", "error getting kubernetes deployment", map[string]interface{}{
					"error": err,
				})
				return false, err
			}
			return deployment.Status.ReadyReplicas == args.Replicas, nil
		},
	})
	if len(errs) > 0 {
		return nil, errs
	}

	multilog.Info("kubernetes.deployments.new", "kubernetes deployment ready", map[string]interface{}{
		"name": result.GetObjectMeta().GetName(),
	})

	return deployment, nil
}

func DeleteDeployment(args DeleteDeploymentArgs) error {
	if err := kubernetes.NewNativeClient().AppsV1().Deployments(args.Namespace).Delete(context.Background(), args.Name, metav1.DeleteOptions{
		PropagationPolicy: variables.ToPtr(metav1.DeletePropagationBackground),
	}); err != nil {
		multilog.Error("kubernetes.deployments.delete", "error deleting kubernetes deployment", map[string]interface{}{
			"error": err,
		})
		return err
	}

	multilog.Info("kubernetes.deployments.delete", "deleted kubernetes deployment", map[string]interface{}{
		"name": args.Name,
	})

	return nil

}
