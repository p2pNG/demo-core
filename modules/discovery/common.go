package discovery

import (
	"github.com/labstack/echo/v4"
	"github.com/p2pNG/core"
	"github.com/p2pNG/core/model"
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
