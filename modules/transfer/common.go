package transfer

import (
	"github.com/labstack/echo/v4"
	"github.com/p2pNG/core"
	"github.com/p2pNG/core/model"
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
