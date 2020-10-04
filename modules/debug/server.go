package debug

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

func dumpClientCertificate(c echo.Context) error {
	x := spew.Sdump(c.Request().TLS.PeerCertificates)
	return c.String(200, x)
}
