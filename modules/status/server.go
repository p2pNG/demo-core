package status

import (
	core "git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getNodeInfo(c echo.Context) error {
	node := NodeInfo{Name: utils.GetHostname(), Core: CoreInfo{
		Version:   core.GetVersionTag(),
		BuildName: core.GetBuildName(),
	}}

	return c.JSONPretty(http.StatusOK, &node, "  ")
}
