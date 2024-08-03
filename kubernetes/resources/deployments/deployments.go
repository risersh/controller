package deployments

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/risersh/controller/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/wait"
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
}

type ResourcesArgs struct {
	Requests ResourceArgs `json:"requests"`
	Limits   ResourceArgs `json:"limits"`
}

type ResourceArgs struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

func NewDeployment(args NewDeploymentArgs) (*appsv1.Deployment, error) {
	client, _ := kubernetes.NewNativeClient()
	deploymentsClient := client.AppsV1().Deployments(args.Namespace)
	replicas := args.Replicas

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      args.Name,
			Namespace: args.Namespace,
			Labels:    args.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
					NodeSelector: map[string]string{
						"role": args.NodeSelector,
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

	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	var readyTime time.Duration
	lastTime := time.Now()

	// wait for deployment to be ready:
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	err = wait.PollUntilContextTimeout(ctx, 1*time.Second, 5*time.Minute, true, func(ctx context.Context) (bool, error) {
		d, err := deploymentsClient.Get(ctx, args.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if d.Status.ReadyReplicas == 1 {
			readyTime = time.Since(lastTime)
			return true, nil
		}

		log.Printf("Waiting for deployment %q to be ready.\n", args.Name)

		return false, nil
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Deployment %q ready after %s.\n", args.Name, readyTime)

	return deployment, nil
}
