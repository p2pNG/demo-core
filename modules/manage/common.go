package manage

import (
	"git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"github.com/labstack/echo/v4"
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
