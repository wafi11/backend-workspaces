package k8s

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateIngress - Create Ingress di K8s
func (k *K8sClient) CreateIngress(ctx context.Context, namespace, name, host, serviceName string, servicePort int) error {
	pathTypePrefix := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                "nginx",
				"cert-manager.io/cluster-issuer":             "letsencrypt-prod",
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{host},
					SecretName: fmt.Sprintf("%s-tls", name), // secret untuk SSL cert
				},
			},
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathTypePrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(servicePort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := k.clientset.NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("ingress %s already exists in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to create ingress: %w", err)
	}

	return nil
}

// GetIngress - Get Ingress by name
func (k *K8sClient) GetIngress(ctx context.Context, namespace, name string) (*networkingv1.Ingress, error) {
	ingress, err := k.clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("ingress %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("failed to get ingress: %w", err)
	}

	return ingress, nil
}

// UpdateIngress - Update Ingress
func (k *K8sClient) UpdateIngress(ctx context.Context, ingress *networkingv1.Ingress) error {
	_, err := k.clientset.NetworkingV1().Ingresses(ingress.Namespace).Update(ctx, ingress, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ingress: %w", err)
	}

	return nil
}

// DeleteIngress - Delete Ingress
func (k *K8sClient) DeleteIngress(ctx context.Context, namespace, name string) error {
	err := k.clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("ingress %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to delete ingress: %w", err)
	}

	return nil
}

// GetIngressURL - Get public URL from Ingress
func (k *K8sClient) GetIngressURL(ctx context.Context, namespace, name string) (string, error) {
	ingress, err := k.GetIngress(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	// Check if ingress has rules
	if len(ingress.Spec.Rules) == 0 {
		return "", fmt.Errorf("ingress has no rules")
	}

	// Get host from first rule
	host := ingress.Spec.Rules[0].Host

	// Check if TLS is enabled
	protocol := "http"
	if len(ingress.Spec.TLS) > 0 {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s", protocol, host), nil
}

// ListIngresses - List all Ingresses in namespace
func (k *K8sClient) ListIngresses(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	ingressList, err := k.clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}

	return ingressList.Items, nil
}
