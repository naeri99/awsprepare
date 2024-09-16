package environ



import (
	"context"
	"log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateEnv(clientset *kubernetes.Clientset) {

	services := []*v1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-common",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					"app": "spark",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-22",
						Port:       22,
						TargetPort: intstrFromInt(22),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-service",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					"app":      "spark",
					"component": "master",
				},
				Ports: []v1.ServicePort{
					{
						Name:       "port-7077",
						Port:       7077,
						TargetPort: intstrFromInt(7077),
					},
					{
						Name:       "port-8083",
						Port:       8083,
						TargetPort: intstrFromInt(8083),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-master-loadbalancer",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: "LoadBalancer",
				Selector: map[string]string{
					"app":      "spark",
					"component": "master",
				},
				Ports: []v1.ServicePort{
					{
						Protocol:   v1.ProtocolTCP,
						Port:       8083,
						TargetPort: intstrFromInt(8083),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "spark-controller-load",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: "LoadBalancer",
				Selector: map[string]string{
					"app":      "spark",
					"component": "master",
					"roleman":  "ui",
				},
				Ports: []v1.ServicePort{
					{
						Protocol:   v1.ProtocolTCP,
						Port:       8889,
						TargetPort: intstrFromInt(8889),
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
