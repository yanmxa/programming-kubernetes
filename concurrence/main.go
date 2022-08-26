package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

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

	// clientSet, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }
	// configMapClient := clientSet.CoreV1().ConfigMaps("default")
	// finalizerConfigMap := &corev1.ConfigMap{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Namespace: "default",
	// 		Name:      "multi-component",
	// 		Labels: map[string]string{
	// 			"labelKey": "labelValue",
	// 		},
	// 	},
	// 	Data: map[string]string{
	// 		"component1": "0",
	// 		"component2": "0",
	// 	},
	// }
	// _, err = configMapClient.Create(context.Background(), finalizerConfigMap, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	messageChan1 := make(chan int)
	messageChan2 := make(chan int)

	ctx := context.Background()
	newContext := context.WithValue(ctx, "chan1", messageChan1)

	var consumeWaitGroup sync.WaitGroup
	consumeWaitGroup.Add(1)
	go consumer1(newContext, &consumeWaitGroup, config)
	// consumeWaitGroup.Add(2)
	// go consumer2(messageChan2, &consumeWaitGroup, config)

	var producerWaitGroup sync.WaitGroup
	// producerNum := 100
	// for i := 0; i < producerNum; i++ {

	// }
	producerWaitGroup.Add(1)
	go produce(messageChan1, messageChan2, &producerWaitGroup)

	producerWaitGroup.Wait()
	close(messageChan1)
	close(messageChan2)
	consumeWaitGroup.Wait()

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	configMapClient := clientSet.CoreV1().ConfigMaps("default")
	latestConfigMap, err := configMapClient.Get(context.TODO(), "multi-component", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	str1 := latestConfigMap.Data["component1"]
	str2 := latestConfigMap.Data["component2"]
	int1, err := strconv.Atoi(str1)
	if err != nil {
		panic(err)
	}
	int2, err := strconv.Atoi(str2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ConfigMap with component1 %d and component2 %d \n", int1, int2)
}

func consumer1(ctx context.Context, wg *sync.WaitGroup, restConfig *rest.Config) {
	ch := ctx.Value("chan1")
	messageChan, ok := ch.(chan int)
	if !ok {
		panic("can not pass chan with context")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	configMapClient := clientSet.CoreV1().ConfigMaps("default")
	for message := range messageChan {
		configMap, err := configMapClient.Get(context.TODO(), "multi-component", metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		str1 := configMap.Data["component1"]
		int1, err := strconv.Atoi(str1)
		if err != nil {
			panic(err)
		}
		fmt.Printf("consumer %d : %s \n", message, strconv.Itoa(int1+1))
		configMap.Data["component1"] = strconv.Itoa(int1 + 1)
		_, err = configMapClient.Update(context.TODO(), configMap, metav1.UpdateOptions{})
		for err != nil {
			fmt.Printf("%+v\n", err)
			_, err = configMapClient.Update(context.TODO(), configMap, metav1.UpdateOptions{})
		}
	}
	wg.Done()
}

func consumer2(messageChan chan int, wg *sync.WaitGroup, restConfig *rest.Config) {
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	configMapClient := clientSet.CoreV1().ConfigMaps("default")

	for message := range messageChan {
		configMap, err := configMapClient.Get(context.TODO(), "multi-component", metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		str2 := configMap.Data["component2"]
		int2, err := strconv.Atoi(str2)
		if err != nil {
			panic(err)
		}
		fmt.Printf("consumer %d : %s \n", message, strconv.Itoa(int2+1))
		configMap.Data["component2"] = strconv.Itoa(int2 + 1)
		_, err = configMapClient.Update(context.TODO(), configMap, metav1.UpdateOptions{})
		for err != nil {
			fmt.Printf("%+v\n", err)
			_, err = configMapClient.Update(context.TODO(), configMap, metav1.UpdateOptions{})
		}
	}
	wg.Done()
}

func produce(messageChan1 chan int, messageChan2 chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	num := 100
	for num := 100; num > 0; num-- {
		// if num%2 != 0 {
		// 	messageChan1 <- 1
		// } else {
		// 	messageChan2 <- 2
		// }
		messageChan1 <- 1
	}
	num++
}
