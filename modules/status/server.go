package status

import (
	"github.com/labstack/echo/v4"
	core "github.com/p2pNG/core"
	"github.com/p2pNG/core/components/database"
	"github.com/p2pNG/core/model"
	"github.com/p2pNG/core/utils"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func getNodeInfo(c echo.Context) error {
	node := model.NodeInfo{
		Name:      utils.GetHostname(),
		Version:   core.GetVersionTag(),
		BuildName: core.GetBuildName(),
	}
	return c.JSONPretty(http.StatusOK, &node, "  ")
}

func listAvailableSeeds(c echo.Context) error {
	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}
	var data [][]byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("file"))
		return bucket.ForEach(func(k, _ []byte) error {
			data = append(data, k)
			return nil
		})
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 5, Message: "read from database error", Internal: err.Error()})
	}
	return c.JSON(http.StatusOK, data)

}
