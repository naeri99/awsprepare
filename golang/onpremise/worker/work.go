package worker

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

func createWorkerEnv(clientset *kubernetes.Clientset) {

	services := []*v1.Service{
		{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-worker-inner",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					"app":       "spark",
					"component": "worker",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-22",
						Protocol:   v1.ProtocolTCP,
						Port:       22,
						TargetPort: intstr.FromInt(2),
					},
				},
			},
		},
		{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-worker-loadbalancer",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer,
				Selector: map[string]string{
					"app":       "spark",
					"component": "worker",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-8090",
						Protocol:   v1.ProtocolTCP,
						Port:       8090,
						TargetPort: intstr.FromInt(8090),
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

func createSparkWorkerStatefulSet(clientset *kubernetes.Clientset) {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spark-worker",
			Namespace: "default",
			Labels: map[string]string{
				"app":       "spark",
				"component": "worker",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "spark-worker-service",
			Replicas:    int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       "spark",
					"component": "worker",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "spark",
						"component": "worker",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "spark-worker",
							Image: "bumory1987/sparks:workerv2",
							Ports: []v1.ContainerPort{
								{ContainerPort: 8083},
								{ContainerPort: 7077},
							},
							Env: []v1.EnvVar{
								{
									Name:  "SPARK_MASTER_HOST",
									Value: "spark-master-0.spark-master-service.default.svc.cluster.local",
								},
								{
									Name:  "SPARK_MASTER_PORT",
									Value: "7077",
								},
								{
									Name:  "SPARK_MASTER_WEBUI_PORT",
									Value: "8083",
								},
								{
									Name:  "SPARK_WORKER_WEBUI_PORT",
									Value: "8090",
								},
								{
									Name:  "STORAGE",
									Value: "spark-storage-0.spark-storage-service.default.svc.cluster.local",
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse("3Gi"),
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

func deleteWrokerStatefulSet(clientset *kubernetes.Clientset) {
	statefulSetName := "spark-worker"
	namespace := "default"

	err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
                fmt.Printf("StatefulSet %s not found\n", statefulSetName)
        }else{
                fmt.Printf("StatefulSet %s deleted\n", statefulSetName)
        }
}

func deleteWorkerEnvSet(clientset *kubernetes.Clientset, target string, namespace string) {
        
	err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), target, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("service %s not found\n", target)
	}else{
		fmt.Printf("service %s deleted\n", target)
	}
}

func deleteWorkerTotalEnv(clientset *kubernetes.Clientset) {
	envList := []string{"spark-worker-loadbalancer", "spark-worker-inner"}
	for _, single := range envList {
		deleteWorkerEnvSet(clientset, single, "default")
	}

}

func DeletingWorker(clientset *kubernetes.Clientset, preoreder chan interface{}) chan interface{} {
	signal := make(chan interface{})
	go func(po chan interface{}, sp chan interface{}, clientset *kubernetes.Clientset) {
		defer close(preoreder)
		<-po
		deleteWorkerTotalEnv(clientset)
		deleteWrokerStatefulSet(clientset)
		sp <- 1
	}(preoreder, signal, clientset)
	return signal
}





func int32Ptr(i int32) *int32 { return &i }

func boolPtr(b bool) *bool { return &b }

func DeployingWorker(clientset *kubernetes.Clientset, preoreder chan interface{}) chan interface{} {
	signal := make(chan interface{})
	go func(po chan interface{}, sp chan interface{}, clientset *kubernetes.Clientset) {
		defer close(preoreder)
		<-po
		createWorkerEnv(clientset)
		createSparkWorkerStatefulSet(clientset)
		sp <- 1
	}(preoreder, signal, clientset)
	return signal
}

