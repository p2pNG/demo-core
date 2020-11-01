package transfer

import (
	"git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"github.com/labstack/echo/v4"
)

type router struct{}

var info = model.PluginInfo{
	Name: "Transfer", Version: "0.0.0", Prefix: "/transfer",
}

func (router) PluginInfo() *model.PluginInfo {
	return &info
}

func (router) GetRouter(g *echo.Group) {
	g.GET("/seed/:hash", getSeed)
	g.GET("/seed/:hash/:block", downloadFileBlock)

}
func init() {
	core.RegisterRouterPlugin(router{})
}
