package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	// For Pods, Services, etc.
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getAccessToken() (string, string) {
	// Step 1: Describe the Kubernetes cluster
	cmd := exec.Command("gcloud", "container", "clusters", "describe", "cluster-1", "--zone", "asia-northeast1-a", "--format", "json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to describe cluster: %v", err)
	}

	// Parse the JSON output to extract the endpoint if needed
	var clusterInfo map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &clusterInfo); err != nil {
		log.Fatalf("Failed to parse cluster info: %v", err)
	}

	// Extract the endpoint (in case you don't know it)
	endpoint := clusterInfo["endpoint"].(string)

	fmt.Printf("Cluster Endpoint: %s\n", endpoint)

	// Step 2: Get the access token
	tokenCmd := exec.Command("gcloud", "auth", "print-access-token")
	var tokenOut bytes.Buffer
	tokenCmd.Stdout = &tokenOut
	err = tokenCmd.Run()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	token := strings.TrimSpace(tokenOut.String())
	fmt.Printf("Access Token: %s\n", token)

	// Step 3: Access the Kubernetes API
	host := "https://" + endpoint // Use the endpoint from the cluster description
	// url := fmt.Sprintf("%s/api/v1/namespaces/default/pods", host)
	return host, token
}

func main() {
	getAccessToken()
	url, token := getAccessToken()

	config := &rest.Config{
		Host:        url,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // Skip SSL verification, not recommended in production
		},
	}

	// Create a clientset using the manually created config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}
	fmt.Println(clientset)
	// List pods in the "default" namespace
	// pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	log.Fatalf("Failed to list pods: %v", err)
	// }

	// fmt.Printf("There are %d pods in the default namespace\n", len(pods.Items))

	// // Print the names of the pods
	// for _, pod := range pods.Items {
	// 	fmt.Printf("Pod name: %s\n", pod.Name)
	// }

	deploymentsClient := clientset.AppsV1().Deployments("default")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	createdDeployment, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create deployment: %v", err)
	}

	fmt.Printf("Created deployment %q.\n", createdDeployment.GetObjectMeta().GetName())

}

func int32Ptr(i int32) *int32 { return &i }

// 	// Example: List Deployments in the "default" namespace
// 	deploymentsClient := clientset.AppsV1().Deployments("default")
// 	// deploymentList, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to list deployments: %v", err)
// 	// }

// 	// // Print the names of the Deployments
// 	// fmt.Printf("There are %d deployments in the default namespace:\n", len(deploymentList.Items))
// 	// for _, d := range deploymentList.Items {
// 	// 	fmt.Printf("Deployment Name: %s\n", d.Name)
// 	// }

// 	// Example: Create a Deployment (replace this part with actual deployment object)
// 	deployment := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "nginx-deployment",
// 		},
// 		Spec: appsv1.DeploymentSpec{
// 			Replicas: int32Ptr(1),
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: map[string]string{
// 					"app": "nginx",
// 				},
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: map[string]string{
// 						"app": "nginx",
// 					},
// 				},
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{
// 						{
// 							Name:  "nginx",
// 							Image: "nginx:1.14.2",
// 							Ports: []corev1.ContainerPort{
// 								{
// 									ContainerPort: 80,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	createdDeployment, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
// 	if err != nil {
// 		log.Fatalf("Failed to create deployment: %v", err)
// 	}

// 	fmt.Printf("Created deployment %q.\n", createdDeployment.GetObjectMeta().GetName())
// }

// func int32Ptr(i int32) *int32 { return &i }

// func main() {
// 	getAccessToken()
// 	url, token := getAccessToken()

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to create HTTP request: %v", err)
// 	}

// 	// Set the Authorization header
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	// Skip SSL verification (equivalent to curl -k)
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	client := &http.Client{Transport: tr}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Failed to make request to Kubernetes API: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Read the response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Failed to read response body: %v", err)
// 	}

// 	// Print the response
// 	fmt.Println(string(body))
// }
