package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sClient) CreateNamespace(ctx context.Context, name string, labels map[string]string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	_, err := k.clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("namespace %s already exists", name)
		}
		return fmt.Errorf("failed to create namespace %s: %w", name, err)
	}

	return nil
}

func (k *K8sClient) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	namespace, err := k.clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("namespace %s not found", name)
		}
		return nil, fmt.Errorf("failed to get namespace %s: %w", name, err)
	}

	return namespace, nil
}

func (k *K8sClient) DeleteNamespace(ctx context.Context, name string) error {
	err := k.clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("namespace %s not found", name)
		}
		return fmt.Errorf("failed to delete namespace %s: %w", name, err)
	}

	return nil
}

func (k *K8sClient) NamespaceExists(ctx context.Context, name string) (bool, error) {
	_, err := k.clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check namespace %s: %w", name, err)
	}

	return true, nil
}

func (k *K8sClient) ListNamespaces(ctx context.Context, labelSelector string) ([]corev1.Namespace, error) {
	namespaceList, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	return namespaceList.Items, nil
}
