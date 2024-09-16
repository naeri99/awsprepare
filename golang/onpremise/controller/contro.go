package controller

import (
	"context"
	"log"
        "fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func createControllEnv(clientset *kubernetes.Clientset) {

	services := []*v1.Service{
		{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-controller-loadbalancer",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer,
				Selector: map[string]string{
					"app":       "spark",
					"component": "ui",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-8889",
						Protocol:   v1.ProtocolTCP,
						Port:       8889,
						TargetPort: intstr.FromInt(8889),
					},
				},
			},
		},
	}

	for _, service := range services {
		_, err := clientset.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
		if err != nil {
			log.Fatalf("Failed to create service %s: %v", service.Name, err)
		}
		log.Printf("Successfully created service: %s", service.Name)
	}
}

func intstrFromInt(i int) intstr.IntOrString {
	return intstr.IntOrString{IntVal: int32(i)}
}

func createSparkControllStatefulSet(clientset *kubernetes.Clientset) {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spark-controller",
			Namespace: "default",
			Labels: map[string]string{
				"app":       "spark",
				"component": "ui",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "spark-controller-loadbalancer",
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       "spark",
					"component": "ui",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "spark",
						"component": "ui",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "spark-controller",
							Image: "bumory1987/sparks:ui",
							Ports: []v1.ContainerPort{
								{ContainerPort: 8889},
							},
							Env: []v1.EnvVar{
								{
									Name:  "STORAGE",
									Value: "spark-storage-0.spark-storage-service.default.svc.cluster.local",
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse("2Gi"),
									v1.ResourceCPU:    resource.MustParse("2"),
								},
							},
							SecurityContext: &v1.SecurityContext{
								Privileged: boolPtr(true),
							},
						},
					},
				},
			},
		},
	}

	// Create the StatefulSet
	_, err := clientset.AppsV1().StatefulSets("default").Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create StatefulSet: %v", err)
	}
	log.Println("Successfully created StatefulSet: spark-storage")
}

func deleteControlStatefulSet(clientset *kubernetes.Clientset){
	statefulSetName := "spark-controller"
	namespace := "default"

	err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
                fmt.Printf("StatefulSet %s not found\n", statefulSetName)
        }else{
                fmt.Printf("StatefulSet %s deleted\n", statefulSetName)
        }
	
}

func deleteControlEnv(clientset *kubernetes.Clientset) {
	serviceName := "spark-controller-loadbalancer"
	namespace := "default"

	err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
                fmt.Printf("service %s not found\n", serviceName)
        }else{
                fmt.Printf("service %s deleted\n", serviceName)
        }
}



func int32Ptr(i int32) *int32 { return &i }

func boolPtr(b bool) *bool { return &b }

func DeployingController(clientset *kubernetes.Clientset, preoreder chan interface{}) chan interface{} {
	signal := make(chan interface{})
	go func(po chan interface{}, sp chan interface{}, clientset *kubernetes.Clientset) {
		defer close(preoreder)
		<-po
		createControllEnv(clientset)
		createSparkControllStatefulSet(clientset)
		sp <- 1
	}(preoreder, signal, clientset)
	return signal
}

func DeletingController(clientset *kubernetes.Clientset) chan interface{} {
        signal := make(chan interface{})
        go func(sp chan interface{}, clientset *kubernetes.Clientset) {
                deleteControlStatefulSet(clientset)
                deleteControlEnv(clientset)
                sp <- 1
        }(signal, clientset)
        return signal
}


