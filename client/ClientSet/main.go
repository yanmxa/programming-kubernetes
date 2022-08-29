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
	// config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	// if err != nil {
	// 	inClusterConfig, err := rest.InClusterConfig()
	// 	if err != nil {
	// 		log.Fatalln("can't get config")
	// 	}
	// 	config = inClusterConfig
	// }

	// clientSet, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	// config.APIPath = "api"
	// config.GroupVersion = &corev1.SchemeGroupVersion
	// config.NegotiatedSerializer = scheme.Codecs
	// restClient, err := rest.RESTClientFor(config)

	// podClient := clientSet.CoreV1().Pods("kube-system")
	// list, err := podClient.List(context.TODO(), metav1.ListOptions{Limit: 500})
	// if err != nil {
	// 	panic(err)
	// }
	// for _, pod := range list.Items {
	// 	fmt.Printf("NameSpace:%v \t Name:%v \t Status:%+v\n", pod.Name, pod.Namespace, pod.Status.Phase)
	// }
	// finalizerConfigMap := &corev1.ConfigMap{
	// 	ObjectMeta: meath1.ObjectMeta{
	// 		Namespace: "default",
	// 		Name:      "test-name",
	// 		Labels: map[string]string{
	// 			"labelKey": "labelValue",
	// 		},
	// 	},
	// 	Data: map[string]string{
	// 		"component1": "",
	// 		"component2": "",
	// 		"component3": "",
	// 	},
	// }
	// dataMap["hello"] = strings.Join([]string{"nihao", "ni", "hao"}, "|")
	// configMapClient := clientSet.CoreV1().ConfigMaps("default")
	// _, err = configMapClient.Create(context.Background(), finalizerConfigMap, meath1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = configMapClient.Get(context.Background(), "test-name", meath1.GetOptions{})
	// if err != nil {
	// 	panic(err)
	// }

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
	fmt.Println(getConfigMap())
}

func getConfigMap() string {
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

	configMapClient := clientSet.CoreV1().ConfigMaps("default")
	configmap, err := configMapClient.Get(context.Background(), "cm", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("get error message %s", err.Error())
	}
	name := configmap.GetName()
	return name
}
