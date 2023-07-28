package main

import (
	"fmt"
	"github.com/Binozo/GoDrop/pkg/air"
	"github.com/Binozo/GoDrop/pkg/airdrop"
	"github.com/Binozo/GoDrop/pkg/owl"
	"time"
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

	// Now we can start discovering
	err = airdrop.StartDiscover(func(receiver air.Receiver) {
		fmt.Println("Receiver found:", receiver.ReceiverName)
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Hour)
}
