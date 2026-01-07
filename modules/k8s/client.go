package k8s

import (
	"context"
	"fmt"

	"github.com/wafi11/backend-workspaces/pkg/k8sclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

// NewK8sClient - Create new K8s client wrapper
func NewK8sClient() (*K8sClient, error) {
	clientset, err := k8sclient.GetClientSet()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s clientset: %w", err)
	}

	return &K8sClient{
		clientset: clientset,
	}, nil
}

// GetClientset - Get underlying clientset
func (k *K8sClient) GetClientset() *kubernetes.Clientset {
	return k.clientset
}

// HealthCheck - Check if K8s cluster is accessible
func (k *K8sClient) HealthCheck(ctx context.Context) error {
	_, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return fmt.Errorf("k8s cluster health check failed: %w", err)
	}
	return nil
}
