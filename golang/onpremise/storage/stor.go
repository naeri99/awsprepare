package storage

import (
	"context"
	"log"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/api/resource" 

)

func createenv(clientset *kubernetes.Clientset) {

	services := []*v1.Service{
		{

			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-storage-service",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "None", // Headless Service
				Selector: map[string]string{
					"app":       "spark",
					"component": "storage",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "netbios-ns",
						Protocol:   v1.ProtocolUDP,
						Port:       137,
						TargetPort: intstr.FromInt(137),
					},
					{
						Name:       "netbios-dgm",
						Protocol:   v1.ProtocolUDP,
						Port:       138,
						TargetPort: intstr.FromInt(138),
					},
					{
						Name:       "netbios-ssn",
						Protocol:   v1.ProtocolTCP,
						Port:       139,
						TargetPort: intstr.FromInt(139),
					},
					{
						Name:       "microsoft-ds",
						Protocol:   v1.ProtocolTCP,
						Port:       445,
						TargetPort: intstr.FromInt(445),
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

func createSparkStorageStatefulSet(clientset *kubernetes.Clientset) {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spark-storage",
			Namespace: "default",
			Labels: map[string]string{
				"app":      "spark",
				"component": "storage",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "spark-storage-service",
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      "spark",
					"component": "storage",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "spark",
						"component": "storage",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "spark-storage",
							Image: "bumory1987/storage:spstorage",
							Ports: []v1.ContainerPort{
								{ContainerPort: 137},
								{ContainerPort: 138},
								{ContainerPort: 139},
								{ContainerPort: 445},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse("2Gi"),
									v1.ResourceCPU:    resource.MustParse("1"),
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


func deleteStorageStatefulSet(clientset *kubernetes.Clientset) {
        statefulSetName := "spark-storage"
        namespace := "default"

        err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
        if err != nil {
                fmt.Printf("StatefulSet %s not found\n", statefulSetName)
        }else{
                fmt.Printf("StatefulSet %s deleted\n", statefulSetName)
        }
}

func deleteStorageEnvSet(clientset *kubernetes.Clientset, target string, namespace string) {
        
        err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), target, metav1.DeleteOptions{})
        if err != nil {
                fmt.Printf("service %s not found\n", target)
        }else{
                fmt.Printf("service %s deleted\n", target)
        }
}

func deleteStorageTotalEnv(clientset *kubernetes.Clientset) {
        envList := []string{"spark-storage-service"}
        for _, single := range envList {
                deleteStorageEnvSet(clientset, single, "default")
        }

}

func DeletingStorage(clientset *kubernetes.Clientset, preoreder chan interface{}) chan interface{} {
        signal := make(chan interface{})
        go func(po chan interface{}, sp chan interface{}, clientset *kubernetes.Clientset) {
                defer close(preoreder)
                <-po
                deleteStorageStatefulSet(clientset)
                deleteStorageTotalEnv(clientset)
                sp <- 1
        }(preoreder, signal, clientset)
        return signal
}







func int32Ptr(i int32) *int32 { return &i }

func boolPtr(b bool) *bool { return &b }


func Deploying(clientset *kubernetes.Clientset) chan interface{} {
	signal := make(chan interface{})
	go func(sp chan interface{}, clientset *kubernetes.Clientset) {
		createenv(clientset)
		createSparkStorageStatefulSet(clientset)
		signal <- 1
	}(signal, clientset)
	return signal
}

