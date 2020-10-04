package debug

import "github.com/labstack/echo/v4"

func GetRouter(e *echo.Group) {
	g := e.Group("/debug")
	g.GET("/client-cert", dumpClientCertificate)
}
