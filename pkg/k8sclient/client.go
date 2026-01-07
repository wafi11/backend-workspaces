package k8sclient

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ClientSet *kubernetes.Clientset
)

func InitK8sClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first (jika app running di dalam K8s cluster)
	config, err = rest.InClusterConfig()
	if err != nil {
		// If not in cluster, try kubeconfig file
		config, err = buildConfigFromFlags()
		if err != nil {
			return nil, fmt.Errorf("failed to build k8s config: %w", err)
		}
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s clientset: %w", err)
	}

	ClientSet = clientset
	return clientset, nil
}

// buildConfigFromFlags - Build config from kubeconfig file
func buildConfigFromFlags() (*rest.Config, error) {
	// Check KUBECONFIG env var
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		// Default to ~/.kube/config
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	// Build config from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig from %s: %w", kubeconfigPath, err)
	}

	return config, nil
}

// GetClientSet - Get initialized clientset
func GetClientSet() (*kubernetes.Clientset, error) {
	if ClientSet == nil {
		return nil, fmt.Errorf("k8s client not initialized, call InitK8sClient() first")
	}
	return ClientSet, nil
}
