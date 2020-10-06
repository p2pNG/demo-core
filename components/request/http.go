package request

import (
	"crypto/tls"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"net/http"
)

func GetDefaultHttpClient() (client *http.Client, err error) {
	if defaultHttpClient == nil {
		err = createDefaultHttpClient()
	}
	client = defaultHttpClient
	return
}

var defaultHttpClient *http.Client

func createDefaultHttpClient() (err error) {
	_, err = certificate.GetCert("client", utils.GetHostname()+" Client")
	if err != nil {
		return
	}
	cert, err := tls.LoadX509KeyPair(certificate.GetCertFilename("client"), certificate.GetCertKeyFilename("client"))
	if err != nil {
		return
	}
	defaultHttpClient = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		},
	}}
	return
}

func DefaultHttpGet() {

}
