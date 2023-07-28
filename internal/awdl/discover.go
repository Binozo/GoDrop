package awdl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Binozo/GoDrop/internal/certificate"
	"github.com/Binozo/GoDrop/pkg/air"
	"github.com/Binozo/GoDrop/pkg/owl"
	"howett.net/plist"
	"io"
	"net"
	"net/http"
	"time"
)

func sendDiscoverRequest(deviceAddress net.IP, port int) (air.Receiver, error) {
	transport, err := certificate.GetTLSTransport()
	if err != nil {
		return air.Receiver{}, err
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	var plistData bytes.Buffer
	encoder := plist.NewBinaryEncoder(&plistData)
	encoder.Encode(map[string][]byte{})

	url := fmt.Sprintf("https://[%s%%25%s]:%d/Discover", deviceAddress.String(), owl.OwlInterface, port)
	request, err := http.NewRequest("POST", url, &plistData)
	if err != nil {
		return air.Receiver{}, err
	}

	request.Header.Add("Content-Type", "application/octet-stream")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("User-Agent", "AirDrop/1.0")
	request.Header.Add("Accept-Language", "en-us")
	request.Header.Add("Accept-Encoding", "br, gzip, deflate")
	request.Close = true

	resp, err := client.Do(request)
	if err != nil {
		return air.Receiver{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return air.Receiver{}, errors.New("invalid status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return air.Receiver{}, err
	}

	decoder := plist.NewDecoder(bytes.NewReader(body))
	var receiverPlistData map[string]interface{}
	decoder.Decode(&receiverPlistData)

	receiverComputerName := receiverPlistData["ReceiverComputerName"].(string) // that's the only useful thing we get/need

	receiver := air.Receiver{
		ReceiverName: receiverComputerName,
		IP:           deviceAddress,
		Port:         port,
		Timestamp:    time.Now(),
	}
	return receiver, nil
}
