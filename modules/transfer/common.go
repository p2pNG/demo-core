package transfer

import "github.com/labstack/echo/v4"

func GetRouter(e *echo.Group) {
	g := e.Group("/transfer")
	g.GET("/seed/:hash", getSeed)
	g.GET("/seed/:hash/:block", downloadFileBlock)
}
