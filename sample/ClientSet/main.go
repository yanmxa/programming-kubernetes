package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("can't get config")
		}
		config = inClusterConfig
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// config.APIPath = "api"
	// config.GroupVersion = &corev1.SchemeGroupVersion
	// config.NegotiatedSerializer = scheme.Codecs
	// restClient, err := rest.RESTClientFor(config)

	podClient := clientSet.CoreV1().Pods("kube-system")
	list, err := podClient.List(context.TODO(), metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}
	for _, pod := range list.Items {
		fmt.Printf("NameSpace:%v \t Name:%v \t Status:%+v\n", pod.Name, pod.Namespace, pod.Status.Phase)
	}

	// result := &corev1.PodList{}
	// err = restClient.Get().
	// 	Namespace("kube-system").
	// 	Resource("pods").
	// 	VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).
	// 	Do(context.TODO()).
	// 	Into(result)

	// if err != nil {
	// 	panic(err)
	// }
	// for _, d := range result.Items {
	// 	fmt.Printf("Namespace: %v \t Name: %v \t Status: %+v \n", d.Namespace, d.Name, d.Status.Phase)
	// }

}
