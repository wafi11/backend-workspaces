package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateSecret - Create Secret di K8s
func (k *K8sClient) CreateSecret(
	ctx context.Context,
	namespace string,
	secretName string,
	secretData map[string]string, // Plain text data, akan di-encode otomatis
) error {
	// Convert plain text to base64
	encodedData := make(map[string][]byte)
	for key, value := range secretData {
		encodedData[key] = []byte(value)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque, // Type untuk generic secret
		Data: encodedData,             // Data dalam format []byte
	}

	_, err := k.clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("secret %s already exists in namespace %s", secretName, namespace)
		}
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

// CreateSecretFromStringData - Alternative: pakai StringData (lebih simple)
func (k *K8sClient) CreateSecretFromStringData(
	ctx context.Context,
	namespace string,
	secretName string,
	secretData map[string]string, // Plain text, K8s akan auto-encode
) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Type:       corev1.SecretTypeOpaque,
		StringData: secretData, // StringData akan di-convert ke Data otomatis
	}

	_, err := k.clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("secret %s already exists in namespace %s", secretName, namespace)
		}
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

// GetSecret - Get Secret by name
func (k *K8sClient) GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	secret, err := k.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("secret %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return secret, nil
}

// GetSecretData - Get decoded secret data
func (k *K8sClient) GetSecretData(ctx context.Context, namespace, name string) (map[string]string, error) {
	secret, err := k.GetSecret(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	// Decode base64 data
	decodedData := make(map[string]string)
	for key, value := range secret.Data {
		decodedData[key] = string(value) // []byte to string
	}

	return decodedData, nil
}

// UpdateSecret - Update Secret
func (k *K8sClient) UpdateSecret(
	ctx context.Context,
	namespace string,
	secretName string,
	secretData map[string]string,
) error {
	// Get existing secret
	secret, err := k.GetSecret(ctx, namespace, secretName)
	if err != nil {
		return err
	}

	// Update StringData
	secret.StringData = secretData

	_, err = k.clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	return nil
}

// DeleteSecret - Delete Secret
func (k *K8sClient) DeleteSecret(ctx context.Context, namespace, name string) error {
	err := k.clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("secret %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// SecretExists - Check if Secret exists
func (k *K8sClient) SecretExists(ctx context.Context, namespace, name string) (bool, error) {
	_, err := k.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check secret: %w", err)
	}

	return true, nil
}
