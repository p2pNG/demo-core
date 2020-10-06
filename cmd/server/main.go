package main

import (
	"crypto/tls"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/components/database"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/debug"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/manage"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/status"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/transfer"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	_, err := database.GetDBEngine()
	if err != nil {
		utils.Log().Fatal("init database error", zap.Error(err))
	}
	defer database.CloseDBEngine()

	go StartHttpServer()
	go func() {
		_, err = discovery.LocalBroadcast(8443)
		if err != nil {
			utils.Log().Fatal("init mDNS service error", zap.Error(err))
		}
	}()
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	<-sign

}

func StartHttpServer() {
	_, _ = certificate.GetCert("server", utils.GetHostname()+" Server")
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Gzip()) //, middleware.Recover())
	api := e.Group("")
	status.GetRouter(api)
	debug.GetRouter(api)
	manage.GetRouter(api)
	transfer.GetRouter(api)
	_ = ReallyStartTlsServer(e, ":8443")
}

func ReallyStartTlsServer(e *echo.Echo, address string) (err error) {
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
