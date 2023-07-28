package airdrop

import "errors"

var notRootError = errors.New("please run this program as root")
var failedCertificateGenerationError = errors.New("failed to generate certificate")
var bleInitFailedError = errors.New("ble init failed")
var failedOWLSetupError = errors.New("failed to setup owl")
var failedStartingBLEAdvertising = errors.New("failed to start BLE advertising")
var failedStartingWebServer = errors.New("failed to start webserver")
