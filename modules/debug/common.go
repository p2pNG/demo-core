package debug

import (
	"git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"github.com/labstack/echo/v4"
)

type router struct{}

var info = model.PluginInfo{
	Name: "Debug", Version: "0.0.0", Prefix: "/debug",
}

func (router) PluginInfo() *model.PluginInfo {
	return &info
}

func (router) GetRouter(g *echo.Group) {
	g.GET("/client-cert", dumpClientCertificate)
}
func init() {
	core.RegisterRouterPlugin(router{})
}
