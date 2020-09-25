package main

import (
	"crypto/tls"
	"fmt"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
)

func main() {
	log.EnableColor()
	log.Info("Hello")
	_, _ = certificate.GetCert("client", "Client Certificate")
	cert, err := tls.LoadX509KeyPair(certificate.GetCertFilename("client"), certificate.GetCertKeyFilename("client"))
	if err != nil {
		log.Error(err)
	}
	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		},
	}}
	resp, _ := client.Get("https://127.0.0.1:8443/")
	text, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(text))
}
