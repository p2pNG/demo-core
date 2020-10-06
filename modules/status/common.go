package status

import "github.com/labstack/echo/v4"

func GetRouter(e *echo.Group) {
	g := e.Group("/status")
	g.GET("/info", getNodeInfo)
	g.GET("/seeds", listAvailableSeeds)
}
