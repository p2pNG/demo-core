package main

import (
	"crypto/tls"
	"fmt"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
)

func main() {
	utils.Log().Info("Hello")
	clients, err := discovery.LocalScan()
	if err != nil {
		utils.Log().Error("scan local peer failed", zap.Error(err))
	}

	_, _ = certificate.GetCert("client", "Client Certificate")
	cert, err := tls.LoadX509KeyPair(certificate.GetCertFilename("client"), certificate.GetCertKeyFilename("client"))
	if err != nil {
		utils.Log().Error("generate local certificate failed", zap.Error(err))
	}
	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		},
	}}
	for cIdx := range clients {
		c := clients[cIdx]
		addr := net.TCPAddr{IP: c.Addr, Port: c.Port}
		endpoint := fmt.Sprintf("https://%s/", addr.String())
		utils.Log().Info("Connecting to " + endpoint)
		resp, _ := client.Get(endpoint)
		text, _ := ioutil.ReadAll(resp.Body)
		fmt.Print(string(text))
	}

}
