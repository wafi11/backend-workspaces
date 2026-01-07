package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateService - Create Service di K8s
func (k *K8sClient) CreateService(
	ctx context.Context,
	namespace string,
	serviceName string,
	appName string,
	port int,
	targetPort int,
) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName, // nama service
			Namespace: namespace,   // namespace
			Labels: map[string]string{
				"app": appName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": appName, // selector untuk connect ke pods dengan label app=appName
			},
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,         // Bukan "http", tapi TCP!
					Port:       int32(port),                // Port service (80)
					TargetPort: intstr.FromInt(targetPort), // Port di container (8080)
				},
			},
		},
	}

	_, err := k.clientset.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("service %s already exists in namespace %s", serviceName, namespace)
		}
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}

// GetService - Get Service by name
func (k *K8sClient) GetService(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	service, err := k.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("service %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return service, nil
}

// DeleteService - Delete Service
func (k *K8sClient) DeleteService(ctx context.Context, namespace, name string) error {
	err := k.clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("service %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

// GetServiceEndpoint - Get service endpoint (ClusterIP:Port)
func (k *K8sClient) GetServiceEndpoint(ctx context.Context, namespace, name string) (string, error) {
	service, err := k.GetService(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	if service.Spec.ClusterIP == "" || service.Spec.ClusterIP == "None" {
		return "", fmt.Errorf("service has no ClusterIP")
	}

	if len(service.Spec.Ports) == 0 {
		return "", fmt.Errorf("service has no ports")
	}

	// Return format: ClusterIP:Port
	return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, service.Spec.Ports[0].Port), nil
}
