package k8s

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentConfig - Configuration untuk create deployment
type DeploymentConfig struct {
	Name          string
	Namespace     string
	AppName       string // Label app
	Image         string // Docker image
	Replicas      int32  // Jumlah pods
	ContainerPort int32  // Port aplikasi di container

	// Resource limits
	CPURequest    string // e.g., "100m", "500m"
	CPULimit      string // e.g., "1000m", "2000m"
	MemoryRequest string // e.g., "128Mi", "256Mi"
	MemoryLimit   string // e.g., "512Mi", "1Gi"

	// Environment variables from ConfigMap
	ConfigMapName string

	// Environment variables from Secret
	SecretName string

	// Custom env vars (optional)
	EnvVars []corev1.EnvVar
}

// CreateDeployment - Create Deployment di K8s
func (k *K8sClient) CreateDeployment(ctx context.Context, config *DeploymentConfig) error {
	// Validate config
	if config.Name == "" || config.Namespace == "" || config.Image == "" {
		return fmt.Errorf("name, namespace, and image are required")
	}

	// Default values
	if config.Replicas == 0 {
		config.Replicas = 1
	}
	if config.ContainerPort == 0 {
		config.ContainerPort = 8080
	}
	if config.CPURequest == "" {
		config.CPURequest = "100m"
	}
	if config.CPULimit == "" {
		config.CPULimit = "500m"
	}
	if config.MemoryRequest == "" {
		config.MemoryRequest = "128Mi"
	}
	if config.MemoryLimit == "" {
		config.MemoryLimit = "512Mi"
	}

	// Build deployment spec
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
			Labels: map[string]string{
				"app": config.AppName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &config.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": config.AppName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": config.AppName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  config.AppName,
							Image: config.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: config.ContainerPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(config.CPURequest),
									corev1.ResourceMemory: resource.MustParse(config.MemoryRequest),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(config.CPULimit),
									corev1.ResourceMemory: resource.MustParse(config.MemoryLimit),
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	// Add EnvFrom (ConfigMap & Secret)
	if config.ConfigMapName != "" || config.SecretName != "" {
		envFrom := []corev1.EnvFromSource{}

		if config.ConfigMapName != "" {
			envFrom = append(envFrom, corev1.EnvFromSource{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: config.ConfigMapName,
					},
				},
			})
		}

		if config.SecretName != "" {
			envFrom = append(envFrom, corev1.EnvFromSource{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: config.SecretName,
					},
				},
			})
		}

		deployment.Spec.Template.Spec.Containers[0].EnvFrom = envFrom
	}

	// Add custom env vars
	if len(config.EnvVars) > 0 {
		deployment.Spec.Template.Spec.Containers[0].Env = config.EnvVars
	}

	// Create deployment
	_, err := k.clientset.AppsV1().Deployments(config.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("deployment %s already exists in namespace %s", config.Name, config.Namespace)
		}
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

// GetDeployment - Get Deployment by name
func (k *K8sClient) GetDeployment(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := k.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("deployment %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

// UpdateDeployment - Update Deployment (e.g., change image, replicas)
func (k *K8sClient) UpdateDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	_, err := k.clientset.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

// UpdateDeploymentImage - Update image (untuk re-deploy dengan image baru)
func (k *K8sClient) UpdateDeploymentImage(ctx context.Context, namespace, name, newImage string) error {
	deployment, err := k.GetDeployment(ctx, namespace, name)
	if err != nil {
		return err
	}

	// Update image
	deployment.Spec.Template.Spec.Containers[0].Image = newImage

	return k.UpdateDeployment(ctx, deployment)
}

// ScaleDeployment - Scale replicas
func (k *K8sClient) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	deployment, err := k.GetDeployment(ctx, namespace, name)
	if err != nil {
		return err
	}

	// Update replicas
	deployment.Spec.Replicas = &replicas

	return k.UpdateDeployment(ctx, deployment)
}

// DeleteDeployment - Delete Deployment
func (k *K8sClient) DeleteDeployment(ctx context.Context, namespace, name string) error {
	err := k.clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("deployment %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	return nil
}

// GetDeploymentStatus - Get status deployment (ready replicas, etc)
func (k *K8sClient) GetDeploymentStatus(ctx context.Context, namespace, name string) (*DeploymentStatus, error) {
	deployment, err := k.GetDeployment(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	status := &DeploymentStatus{
		Name:              deployment.Name,
		Namespace:         deployment.Namespace,
		Replicas:          deployment.Status.Replicas,
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
		UpdatedReplicas:   deployment.Status.UpdatedReplicas,
		IsReady:           deployment.Status.ReadyReplicas == *deployment.Spec.Replicas,
	}

	return status, nil
}

// WaitForDeploymentReady - Wait until deployment is ready
func (k *K8sClient) WaitForDeploymentReady(ctx context.Context, namespace, name string, timeoutSeconds int) error {
	// Implementation using watch API
	// Simplified version: poll setiap 2 detik

	// Note: Untuk production, pakai watch API
	// Ini contoh simplified dengan polling

	return fmt.Errorf("not implemented yet - use kubectl wait or watch API")
}

// ListDeployments - List all deployments in namespace
func (k *K8sClient) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	deploymentList, err := k.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	return deploymentList.Items, nil
}

// DeploymentStatus - Status deployment
type DeploymentStatus struct {
	Name              string
	Namespace         string
	Replicas          int32
	ReadyReplicas     int32
	AvailableReplicas int32
	UpdatedReplicas   int32
	IsReady           bool
}
