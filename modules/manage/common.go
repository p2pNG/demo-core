package manage

import (
	"github.com/labstack/echo/v4"
	"github.com/p2pNG/core"
	"github.com/p2pNG/core/model"
)

type router struct{}

var info = model.PluginInfo{
	Name: "Manage", Version: "0.0.0", Prefix: "/manage",
}

func (router) PluginInfo() *model.PluginInfo {
	return &info
}

func (router) GetRouter(g *echo.Group) {
	g.GET("/add-file", addLocalFile)

}
func init() {
	core.RegisterRouterPlugin(router{})
}
