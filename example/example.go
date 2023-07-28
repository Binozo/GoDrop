package main

import (
	"fmt"
	"github.com/Binozo/GoDrop/pkg/air"
	"github.com/Binozo/GoDrop/pkg/airdrop"
	"github.com/Binozo/GoDrop/pkg/owl"
	"github.com/korylprince/go-cpio-odc"
)

func main() {
	// Configure your GoDrop setup
	// (You can define the Wifi interface here)
	airdrop.Config(air.Config{DebugMode: true, HostName: "GoDrop"})

	// Let's subscribe to OWL-errors
	// Great for debugging and checking if
	// OWL works correctly on your hardware
	//
	// Take a look `owl.SubscribeToLogs()` if you want
	// to listen to logs
	go func() {
		for {
			owlError := <-owl.SubscribeToErrors()
			fmt.Println("An owl error occurred:", owlError.ErrorType)
			fmt.Println("Logs:", owlError.Logs)
		}
	}()

	// Init everything required (wifi, owl, ble...)
	err := airdrop.Init()
	if err != nil {
		panic(err)
	}

	// Here you can accept or decline the sending request
	// If you don't define this all requests will be declined
	airdrop.OnAsk(func(request air.Request) bool {
		fmt.Println("Incoming request from", request.SenderModelName)
		fmt.Println("Wants to send", len(request.Files), "files")
		return true
	})

	// Here you will be able to access the received files
	// Note that OnAsk will be called first, after that
	// you will get the files here
	airdrop.OnFiles(func(request air.Request, files []*cpio.File) {
		fmt.Println("Got files from", request.SenderComputerName)
		for _, file := range files {
			fmt.Println("Filename:", file.Name())
		}
	})

	err = airdrop.StartReceiving()
	if err != nil {
		panic(err)
	}
}
