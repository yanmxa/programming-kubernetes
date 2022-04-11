package main

import (
	"fmt"
	"log"

	// "programming-kubernetes/pkg"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	// 1. config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("can't get config")
		}
		config = inClusterConfig
	}

	// 2. client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can't create client")
	}

	// 3. get informer
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("default"))
	podInformer := factory.Core().V1().Pods().Informer()
	// factory := informers.NewSharedInformerFactory(clientset, 0)
	// serviceInformer := factory.Core().V1().Services()
	// ingressInformer := factory.Networking().V1().Ingresses()

	// add workQueue
	rateLimitQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Controller")

	// 4. add event handler
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Add Event")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Println("can't get key")
			}
			rateLimitQueue.AddRateLimited(key)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Update Event")
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				fmt.Println("can't get key")
			}
			rateLimitQueue.AddRateLimited(key)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete Event")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Println("can't get key")
			}
			rateLimitQueue.AddRateLimited(key)
		},
	})
	// controller := pkg.NewController(clientset, serviceInformer, ingressInformer)

	stopCh := make(chan struct{})

	// 5. informer.Start
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// controller.Run(stopCh)
}
