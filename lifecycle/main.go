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
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV, syscall.SIGALRM)
	go func() {
		for {
			s := <-signalChan
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[ %s ] received signal: %d \n", timeStr, s)
		}
	}()

	for {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[ %s ] container running ", timeStr)
	}
}
