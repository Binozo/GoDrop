package certificate

import (
	"fmt"
	"os"
	"os/exec"
)

const KeyFileName = "key.pem"
const CertificateFileName = "certificate.pem"

// CertificatesPresent checks if the required TLS certificate is present
func CertificatesPresent() bool {
	if _, err := os.Stat(KeyFileName); err != nil {
		return false
	}
	return true
}

// GenerateCertificates creates a new TLS certificate using openssl.
// No Apple Key extraction required.
// Thanks to https://github.com/seemoo-lab/opendrop/blob/cbc7ca46a75bb2da4af9d406732ddf2192578388/opendrop/config.py#L136
func GenerateCertificates(hostname string) error {
	cmd := exec.Command("openssl", "req", "-newkey", "rsa:2048", "-nodes", "-keyout", KeyFileName, "-x509", "-days", "3650", "-out", CertificateFileName, "-subj", fmt.Sprintf("/CN=%s", hostname))
	err := cmd.Run()
	if err != nil {
		// Failed
		return err
	}
	return nil
}
