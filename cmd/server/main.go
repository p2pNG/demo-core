package main

import (
	"crypto/tls"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
)

func main() {
	go StartHttpServer()
	go func() {
		_, _ = discovery.LocalBroadcast(8443)
	}()
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	<-sign

}

func StartHttpServer() {
	_, _ = certificate.GetCert("server", "Server Certificate")
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		x := spew.Sdump(c.Request().TLS.PeerCertificates)
		return c.String(200, x)
	})
	if e.Server.TLSConfig == nil {
		e.Server.TLSConfig = &tls.Config{}
	}

	e.Server.TLSConfig.ClientAuth = tls.RequireAnyClientCert
	e.Server.Addr = ":8443"
	_ = e.Server.ListenAndServeTLS(certificate.GetCertFilename("server"), certificate.GetCertKeyFilename("server"))
}
