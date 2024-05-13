# GoDrop
A Go implementation of [OpenDrop](https://github.com/seemoo-lab/opendrop).

[![Go Reference](https://pkg.go.dev/badge/github.com/Binozo/GoDrop.svg)](https://pkg.go.dev/github.com/Binozo/GoDrop)

> [!WARNING]
> This library only works until macOS Ventura and iOS 15. Apple apparently changed the way their AirDrop protocol works.

Big thanks to [Seemoo 🇩🇪](https://github.com/seemoo-lab) for developing Open-Source implementations for AirDrop.
###### Made In Germany

This project is inspired by [GoOpenDrop](https://github.com/bodaay/GoOpenDrop).

## Features
- 🔓 Works without the need of [extraction](https://github.com/seemoo-lab/airdrop-keychain-extractor) of Apple signed certificates
- ⚡ Easy High-Level Go API
- 💪 Manages everything

#### 🔌 Note: This project is designed to be extended by you
#### ✔️ Works the best on a Raspberry Pi

## Setup
### Prerequisites
This project has been developed and tested on a Raspberry Pi 4B 2GB. <br>
Other platforms may not work as good as on a Raspberry Pi.

You **NEED** a Wifi card with **ACTIVE MONITORING**.
###### (I didn't test Nexamon patched firmwares.)

#### How do I know my Wifi Card supports _Active Monitoring_?
Type `iw list` into your Terminal. <br>
Now check if you can find `Device supports active monitor (which will ACK incoming frames)` 
in the entry of your Wifi Card. 

One-liner:
```bash
iw list | grep -q "Device supports active monitor"; echo $?
```
Does it return `0`? Your Wifi Card supports active monitor mode. <br>
If it returns `1` it does not.

I am using the `ALLNET ALL-WA1200AC` Wifi Card. It is plug-and-play.

### Dependencies
#### Go
First make sure you got Go installed on your machine. <br>
If you want to install Go quickly you can try using this [script](https://github.com/canha/golang-tools-install-script#fast_forward-install).

[This](https://github.com/canha/golang-tools-install-script#fast_forward-install) repository offers you a quick installation script:
```bash
wget -q -O - https://git.io/vQhTU | bash
```

#### Packages
Now install the required packages:
```bash
sudo apt install imagemagick bluez
```

#### OWL
Follow [OWL](https://github.com/seemoo-lab/owl)'s setup [here](https://github.com/seemoo-lab/owl#requirements).

Ideally make a quick test if owl works:
```bash
sudo owl -i <WIFI_INTERFACE> -c 44 -v
```
If everything seems to work hit `CTLR + C`.

#### Go Module
Run
```bash
go get -u github.com/Binozo/GoDrop
```
in your Project.

## Usage
Now you can finally use this library 🥳 <br>

Example:
```go
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
		fmt.Println("Incoming request from", request.SenderComputerName)
		fmt.Println("Wants to send", len(request.Files), "files")
		return true // accepting
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
```
##### (You can find this file here: [example/example.go](https://github.com/Binozo/GoDrop/tree/master/example/example.go))

### TODO's
- [x] Add discovering receiver functionality
- [ ] Improve discovering [(sometimes the receiver doesn't react)](https://github.com/Binozo/GoDrop/blob/b6088aa7652392c2edd4e20268a87abca74c51df/internal/awdl/service.go#L76)
- [ ] Add uploading feature (first we need to improve discovering)