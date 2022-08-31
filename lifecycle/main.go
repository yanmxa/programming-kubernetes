package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV,
		syscall.SIGALRM)

	fmt.Printf("# [ %s ] main container start running \n", time.Now().Format("2006-01-02 15:04:05"))

	cfg, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("failed to get kube config: %v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Printf("failed to get kube client: %v", err)
		return
	}

	ticker := time.NewTicker(time.Second)
	count := 0
loop:
	for {
		select {
		case sig := <-signalChan:
			fmt.Printf("# [ %s ] receive signal: %s => %d \n", time.Now().Format("2006-01-02 15:04:05"), sig.String(), sig)
			count = 0
			break loop
		case <-ticker.C:
			count++
			fmt.Printf(" #(%d) ", count)
		}
	}

	fmt.Printf("# [ %s ] graceful shudown >>>>>>>>>>>>>>>>>>> \n", time.Now().Format("2006-01-02 15:04:05"))
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		count++
		sa, _ := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), "default", metav1.GetOptions{})
		fmt.Printf(" #(%d) => service account: %s \n", count, sa.Name)
	}

	fmt.Printf("# [ %s ] main container finished \n", time.Now().Format("2006-01-02 15:04:05"))
}
