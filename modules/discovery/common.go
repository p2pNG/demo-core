package discovery

import (
	"git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"github.com/labstack/echo/v4"
)

type router struct{}

var info = model.PluginInfo{
	Name: "Discovery", Version: "0.0.0", Prefix: "/discovery",
}

func (router) PluginInfo() *model.PluginInfo {
	return &info
}

func (router) GetRouter(g *echo.Group) {
	g.POST("/register", registerClient)
	g.GET("/peers", listAvailablePeers)
}
func init() {
	core.RegisterRouterPlugin(router{})
}
