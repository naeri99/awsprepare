package login

import (
	"io/ioutil"
	"log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)


func initialization()  *rest.Config{
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("Failed to read token from %s: %v", tokenPath, err)
	}

	caCertPath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to read CA cert from %s: %v", caCertPath, err)
	}

	kubeAPIServer := "https://kubernetes.default.svc"


	config := &rest.Config{
		Host:        kubeAPIServer,              
		BearerToken: string(token),              
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCert,                     
		},
	}

	return config

}

func GetClient() *kubernetes.Clientset{
	config := initialization()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	return clientset
}

