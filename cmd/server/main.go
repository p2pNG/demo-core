package main

import (
	"crypto/tls"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/debug"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/status"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
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
	_, _ = certificate.GetCert("server", utils.GetHostname()+" Server")
	e := echo.New()
	api := e.Group("")
	status.GetRouter(api)
	debug.GetRouter(api)
	_ = ReallyStartTlsSServer(e, ":8443")
}

func ReallyStartTlsSServer(e *echo.Echo, address string) (err error) {
	s := e.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)
	cert := certificate.GetCertFilename("server")
	key := certificate.GetCertKeyFilename("server")
	if s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(cert, key); err != nil {
		return
	}
	s.TLSConfig.ClientAuth = tls.RequireAnyClientCert

	s.Addr = address
	if !e.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return e.StartServer(e.TLSServer)
}
