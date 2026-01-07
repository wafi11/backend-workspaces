package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wafi11/backend-workspaces/modules/k8s"
	"github.com/wafi11/backend-workspaces/pkg/k8sclient"
	corev1 "k8s.io/api/core/v1"
)

func main() {
	// ============================================
	// 1. Initialize K8s Client
	// ============================================
	log.Println("🚀 Initializing K8s client...")

	_, err := k8sclient.InitK8sClient()
	if err != nil {
		log.Fatalf("❌ Failed to init k8s client: %v", err)
	}

	k8sClient, err := k8s.NewK8sClient()
	if err != nil {
		log.Fatalf("❌ Failed to create k8s client: %v", err)
	}

	ctx := context.Background()

	// Health check
	if err := k8sClient.HealthCheck(ctx); err != nil {
		log.Fatalf("❌ K8s health check failed: %v", err)
	}

	log.Println("✅ K8s client initialized successfully")

	// ============================================
	// 2. Define App Configuration
	// ============================================
	namespace := "user123-mystore"
	appName := "mystore"
	subdomain := "mystore.yourplatform.com" // Ganti dengan domain kamu
	dockerImage := "nginx:latest"           // Ganti dengan image kamu nanti

	log.Printf("📦 Deploying app: %s to namespace: %s\n", appName, namespace)

	// ============================================
	// 3. Create Namespace
	// ============================================
	log.Println("\n📁 Step 1/6: Creating namespace...")

	err = k8sClient.CreateNamespace(ctx, namespace, map[string]string{
		"app":        appName,
		"user":       "user123",
		"managed-by": "your-platform",
	})
	if err != nil {
		log.Printf("⚠️  Namespace creation warning: %v (might already exist)\n", err)
	} else {
		log.Println("✅ Namespace created successfully")
	}

	// ============================================
	// 4. Create ConfigMap (Non-sensitive config)
	// ============================================
	log.Println("\n⚙️  Step 2/6: Creating ConfigMap...")

	configData := map[string]string{
		"APP_NAME":  "My Store",
		"APP_ENV":   "production",
		"LOG_LEVEL": "info",
		"DB_HOST":   "postgres-service",
		"DB_PORT":   "5432",
		"DB_NAME":   "mystore",
		"SMTP_HOST": "smtp.gmail.com",
		"SMTP_PORT": "587",
	}

	err = k8sClient.CreateConfigMap(
		ctx,
		namespace,
		fmt.Sprintf("%s-config", appName),
		configData,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create configmap: %v", err)
	}

	log.Println("✅ ConfigMap created successfully")

	// ============================================
	// 5. Create Secret (Sensitive data)
	// ============================================
	log.Println("\n🔐 Step 3/6: Creating Secret...")

	secretData := map[string]string{
		"STRIPE_API_KEY":        "sk_test_xxxxxxxxxxxxxxxx",
		"STRIPE_WEBHOOK_SECRET": "whsec_xxxxxxxxxxxxxxxx",
		"DATABASE_PASSWORD":     "superSecretPassword123",
		"JWT_SECRET":            "myJwtSecretKey789",
		"SMTP_PASSWORD":         "emailPassword456",
	}

	err = k8sClient.CreateSecretFromStringData(
		ctx,
		namespace,
		fmt.Sprintf("%s-secrets", appName),
		secretData,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create secret: %v", err)
	}

	log.Println("✅ Secret created successfully")

	// ============================================
	// 6. Create Deployment
	// ============================================
	log.Println("\n🚢 Step 4/6: Creating Deployment...")

	deploymentConfig := &k8s.DeploymentConfig{
		Name:          fmt.Sprintf("%s-deployment", appName),
		Namespace:     namespace,
		AppName:       appName,
		Image:         dockerImage,
		Replicas:      2,  // 2 pods untuk high availability
		ContainerPort: 80, // nginx listen di port 80

		// Resource limits
		CPURequest:    "100m",
		CPULimit:      "500m",
		MemoryRequest: "128Mi",
		MemoryLimit:   "512Mi",

		// Inject ConfigMap & Secret
		ConfigMapName: fmt.Sprintf("%s-config", appName),
		SecretName:    fmt.Sprintf("%s-secrets", appName),

		// Custom env vars (optional)
		EnvVars: []corev1.EnvVar{
			{
				Name:  "PORT",
				Value: "80",
			},
		},
	}

	err = k8sClient.CreateDeployment(ctx, deploymentConfig)
	if err != nil {
		log.Fatalf("❌ Failed to create deployment: %v", err)
	}

	log.Println("✅ Deployment created successfully")

	// Wait a bit for pods to start
	log.Println("⏳ Waiting for deployment to be ready...")
	// Kamu bisa implement wait logic disini atau pakai kubectl wait

	// ============================================
	// 7. Create Service
	// ============================================
	log.Println("\n🔌 Step 5/6: Creating Service...")

	err = k8sClient.CreateService(
		ctx,
		namespace,
		fmt.Sprintf("%s-service", appName), // service name
		appName,                            // app label (selector)
		80,                                 // service port
		80,                                 // target port (nginx port)
	)
	if err != nil {
		log.Fatalf("❌ Failed to create service: %v", err)
	}

	log.Println("✅ Service created successfully")

	// ============================================
	// 8. Create Ingress
	// ============================================
	log.Println("\n🌐 Step 6/6: Creating Ingress...")

	err = k8sClient.CreateIngress(
		ctx,
		namespace,
		fmt.Sprintf("%s-ingress", appName), // ingress name
		subdomain,                          // host (mystore.yourplatform.com)
		fmt.Sprintf("%s-service", appName), // service name
		80,                                 // service port
	)
	if err != nil {
		log.Fatalf("❌ Failed to create ingress: %v", err)
	}

	log.Println("✅ Ingress created successfully")

	// ============================================
	// 9. Get Deployment Status
	// ============================================
	log.Println("\n📊 Getting deployment status...")

	status, err := k8sClient.GetDeploymentStatus(
		ctx,
		namespace,
		fmt.Sprintf("%s-deployment", appName),
	)
	if err != nil {
		log.Printf("⚠️  Failed to get deployment status: %v\n", err)
	} else {
		log.Printf("   Total Replicas: %d\n", status.Replicas)
		log.Printf("   Ready Replicas: %d\n", status.ReadyReplicas)
		log.Printf("   Is Ready: %v\n", status.IsReady)
	}

	// ============================================
	// 10. Get Ingress URL
	// ============================================
	ingressURL, err := k8sClient.GetIngressURL(
		ctx,
		namespace,
		fmt.Sprintf("%s-ingress", appName),
	)
	if err != nil {
		log.Printf("⚠️  Failed to get ingress URL: %v\n", err)
	} else {
		log.Printf("\n🎉 Deployment completed successfully!")
		log.Printf("🌐 Your app is available at: %s\n", ingressURL)
	}

	// ============================================
	// 11. Summary
	// ============================================
	log.Println("\n" + "=======================")
	log.Println("📋 DEPLOYMENT SUMMARY")
	log.Println("=========================")
	log.Printf("Namespace:    %s\n", namespace)
	log.Printf("App Name:     %s\n", appName)
	log.Printf("ConfigMap:    %s-config\n", appName)
	log.Printf("Secret:       %s-secrets\n", appName)
	log.Printf("Deployment:   %s-deployment\n", appName)
	log.Printf("Service:      %s-service\n", appName)
	log.Printf("Ingress:      %s-ingress\n", appName)
	log.Printf("Public URL:   %s\n", ingressURL)
	log.Println("=======================")

	// ============================================
	// 12. Kubectl Commands (untuk verify)
	// ============================================
	log.Println("\n💡 Verify deployment dengan commands berikut:")
	log.Printf("   kubectl get all -n %s\n", namespace)
	log.Printf("   kubectl get pods -n %s\n", namespace)
	log.Printf("   kubectl get ingress -n %s\n", namespace)
	log.Printf("   kubectl logs -f deployment/%s-deployment -n %s\n", appName, namespace)
}
