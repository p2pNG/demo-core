package main

import (
	"crypto/tls"
	"errors"
	"git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/components/database"

	_ "git.ixarea.com/p2pNG/p2pNG-core/cmd/plugins"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
)

const ListenPort = 6444

func main() {
	err := os.MkdirAll(utils.AppDataDir(), 0755)
	if err != nil {
		utils.Log().Fatal("init data path error", zap.Error(err))
	}

	_, err = database.GetDBEngine()
	if err != nil {
		utils.Log().Fatal("init database error", zap.Error(err))
	}
	defer database.CloseDBEngine()

	go func() {
		err := StartHttpServer()
		if err != nil {
			utils.Log().Fatal("init http  service error", zap.Error(err))
		}
	}()
	go func() {
		_, err = discovery.LocalBroadcast(ListenPort)
		if err != nil {
			utils.Log().Fatal("init mDNS service error", zap.Error(err))
		}
	}()
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	<-sign

}

func StartHttpServer() error {
	_, _ = certificate.GetCert("server", utils.GetHostname()+" Server")
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Gzip(), middleware.Recover())
	api := e.Group("")

	plugins := []string{"Debug", "Discovery", "Status", "Manage", "Transfer"}
	for nameIdx := range plugins {
		name := plugins[nameIdx]
		p, ok := core.GetRouterPlugin(name)

		if !ok {
			return errors.New("unregistered plugin: " + name)

		}
		info := p.PluginInfo()
		utils.Log().Info("loading router plugin", zap.String("plugin", name))
		p.GetRouter(api.Group(info.Prefix))
	}

	_ = ReallyStartTlsServer(e, ":"+strconv.Itoa(ListenPort))
	return nil
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
