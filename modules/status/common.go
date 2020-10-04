package status

import "github.com/labstack/echo/v4"

type NodeInfo struct {
	Name string
	Core CoreInfo
}
type CoreInfo struct {
	Version   string
	BuildName string
}

func GetRouter(e *echo.Group) {
	g := e.Group("/status")
	g.GET("/info", getNodeInfo)
}
