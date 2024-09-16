package master

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

func createMasterEnv(clientset *kubernetes.Clientset) {

	services := []*v1.Service{
		{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-master-service",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					"app":       "spark",
					"component": "master",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-7077",
						Protocol:   v1.ProtocolTCP,
						Port:       7077,
						TargetPort: intstr.FromInt(7077),
					},
					{
						Name:       "port-8083",
						Protocol:   v1.ProtocolTCP,
						Port:       8083,
						TargetPort: intstr.FromInt(8083),
					},
				},
			},
		},
		{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-master-external",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer,
				Selector: map[string]string{
					"app":       "spark",
					"component": "master",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-7077",
						Protocol:   v1.ProtocolTCP,
						Port:       8083,
						TargetPort: intstr.FromInt(8083),
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

func createSparkMasterStatefulSet(clientset *kubernetes.Clientset) {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spark-master",
			Namespace: "default",
			Labels: map[string]string{
				"app":       "spark",
				"component": "master",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "spark-master-service",
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       "spark",
					"component": "master",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "spark",
						"component": "master",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "spark-master",
							Image: "bumory1987/sparks:masterv2",
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


func deleteMasterStatefulSet(clientset *kubernetes.Clientset) {
        statefulSetName := "spark-master"
        namespace := "default"

        err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
        if err != nil {
                fmt.Printf("StatefulSet %s not found\n", statefulSetName)
        }else{
                fmt.Printf("StatefulSet %s deleted\n", statefulSetName)
        }
}

func deleteMasterEnvSet(clientset *kubernetes.Clientset, target string, namespace string) {
        
        err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), target, metav1.DeleteOptions{})
        if err != nil {
                fmt.Printf("service %s not found\n", target)
        }else{
                fmt.Printf("service %s deleted\n", target)
        }
}

func deleteMasterTotalEnv(clientset *kubernetes.Clientset) {
        envList := []string{"spark-master-external", "spark-master-service"}
        for _, single := range envList {
                deleteMasterEnvSet(clientset, single, "default")
        }

}

func DeletingMaster(clientset *kubernetes.Clientset, preoreder chan interface{}) chan interface{} {
        signal := make(chan interface{})
        go func(po chan interface{}, sp chan interface{}, clientset *kubernetes.Clientset) {
                defer close(preoreder)
                <-po
		deleteMasterStatefulSet(clientset)
                deleteMasterTotalEnv(clientset)
                sp <- 1
        }(preoreder, signal, clientset)
        return signal
}


func int32Ptr(i int32) *int32 { return &i }

func boolPtr(b bool) *bool { return &b }

func DeployingMaster(clientset *kubernetes.Clientset, preoreder chan interface{}) chan interface{} {
	signal := make(chan interface{})
	go func(po chan interface{}, sp chan interface{}, clientset *kubernetes.Clientset) {
		defer close(preoreder)
		<-po
		createMasterEnv(clientset)
		createSparkMasterStatefulSet(clientset)
		signal <- 1
	}(preoreder, signal, clientset)
	return signal
}

