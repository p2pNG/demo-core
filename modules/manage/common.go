package manage

import "github.com/labstack/echo/v4"

func GetRouter(e *echo.Group) {
	g := e.Group("/manage")
	g.GET("/add-file", addLocalFile)
}
