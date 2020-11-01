package status

import (
	"git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"github.com/labstack/echo/v4"
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
