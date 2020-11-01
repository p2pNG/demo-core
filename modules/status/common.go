package status

import (
	"github.com/labstack/echo/v4"
	"github.com/p2pNG/core"
	"github.com/p2pNG/core/model"
)

type router struct{}

var info = model.PluginInfo{
	Name: "Status", Version: "0.0.0", Prefix: "/status",
}

func (router) PluginInfo() *model.PluginInfo {
	return &info
}

func (router) GetRouter(g *echo.Group) {
	g.GET("/info", getNodeInfo)
	g.GET("/seeds", listAvailableSeeds)

}
func init() {
	core.RegisterRouterPlugin(router{})
}
