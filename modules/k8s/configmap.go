package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateConfigMap - Create ConfigMap di K8s
func (k *K8sClient) CreateConfigMap(ctx context.Context, namespace, name string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	_, err := k.clientset.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("configmap %s already exists in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to create configmap: %w", err)
	}

	return nil
}

// GetConfigMap - Get ConfigMap by name
func (k *K8sClient) GetConfigMap(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	configMap, err := k.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("configmap %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}

	return configMap, nil
}

func (k *K8sClient) UpdateConfigMap(ctx context.Context, namespace, name string, data map[string]string) error {
	configMap, err := k.GetConfigMap(ctx, namespace, name)
	if err != nil {
		return err
	}

	configMap.Data = data

	_, err = k.clientset.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update configmap: %w", err)
	}

	return nil
}

func (k *K8sClient) DeleteConfigMap(ctx context.Context, namespace, name string) error {
	err := k.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("configmap %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to delete configmap: %w", err)
	}

	return nil
}

func (k *K8sClient) ConfigMapExists(ctx context.Context, namespace, name string) (bool, error) {
	_, err := k.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check configmap: %w", err)
	}

	return true, nil
}
