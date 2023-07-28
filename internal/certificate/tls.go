package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

func GetTLSTransport() (*http.Transport, error) {
	cert, err := tls.LoadX509KeyPair(CertificateFileName, KeyFileName)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(appleRootCert)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		},
		TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{}, //disable http2
	}
	return tr, nil
}
