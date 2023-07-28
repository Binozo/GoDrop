package webserver

import (
	"fmt"
	"github.com/Binozo/GoDrop/internal/awdl"
	"github.com/Binozo/GoDrop/internal/certificate"
	"github.com/Binozo/GoDrop/pkg/owl"
	"github.com/gorilla/mux"
	"net/http"
)

// Start the webserver for accepting AirDrop files
func Start() error {
	r := setup()
	address := fmt.Sprintf("[%s%%%s]:%d", owl.GetOwlInterfaceAddr(), owl.OwlInterface, awdl.AirDropPort)
	//log.Debug().Msgf("Serving at %s", address)

	return http.ListenAndServeTLS(address, certificate.CertificateFileName, certificate.KeyFileName, r)
}

func setup() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", rootHandler).Methods("HEAD")
	r.HandleFunc("/Discover", discoverHandler).Methods("POST")
	r.HandleFunc("/Ask", askHandler).Methods("POST")
	r.HandleFunc("/Upload", uploadHandler).Methods("POST")
	return r
}
