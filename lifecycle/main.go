package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV, syscall.SIGALRM, syscall.SIGKILL)

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("# [ %s ] main container start running \n", timeStr)

	ticker := time.NewTicker(time.Second)
	count := 0
	// loop:
	for {
		select {
		case sig := <-signalChan:
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("# [ %s ] receive signal: %s => %d \n", timeStr, sig.String(), sig)
			count = 0
			// break loop
		case <-ticker.C:
			count++
			fmt.Printf(" #(%d) ", count)
		}
	}
	// time.Sleep(5 * time.Second)
	// timeStr = time.Now().Format("2006-01-02 15:04:05")
	// fmt.Printf("# [ %s ] main container finished \n", timeStr)
}

// journalctl -f -u kubelet
// kubectl create configmap cm --from-literal=special.how=very --from-literal=special.type=charm
// func getConfigMap() string {
// 	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
// 	if err != nil {
// 		inClusterConfig, err := rest.InClusterConfig()
// 		if err != nil {
// 			log.Fatalln("can't get config")
// 		}
// 		config = inClusterConfig
// 	}

// 	clientSet, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	configMapClient := clientSet.CoreV1().ConfigMaps("default")
// 	configmap, err := configMapClient.Get(context.Background(), "cm", metav1.GetOptions{})
// 	if err != nil {
// 		fmt.Printf("get error message %s", err.Error())
// 	}
// 	name := configmap.GetName()
// 	return name
// }
