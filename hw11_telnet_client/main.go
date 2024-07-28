package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	i := 1
	if len(os.Args) > 3 {
		i = 2
	}
	if len(os.Args) < 3 {
		panic("Wrong args")
	}
	timeoutString := flag.String("timeout", "10s", "usage string")
	flag.Parse()
	host := os.Args[i]
	port := os.Args[i+1]

	timeout, err := time.ParseDuration(*timeoutString)
	if err != nil {
		panic(err.Error())
	}

	client := NewTelnetClient(fmt.Sprintf("%s:%s", host, port), timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		fmt.Printf("Connection error: %v", err)
		return
	}
	defer client.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go listenSignals(client, wg)
	go listenFromClient(client, wg)
	go sendToClient(client, wg)
	go forceTimeout(timeout, wg)
	wg.Wait()
}

func listenSignals(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	<-signalChannel
	err := client.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func listenFromClient(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Receive()
	if err != nil {
		client.Close()
	}
}

func sendToClient(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Send()
	if err != nil {
		fmt.Println(err.Error())
	}
	client.Close()
}

func forceTimeout(timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(timeout)
}
