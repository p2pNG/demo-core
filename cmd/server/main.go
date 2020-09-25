package main

import (
	"crypto/tls"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
)

func main() {

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

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	<-sign
}
