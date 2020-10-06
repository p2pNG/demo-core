package manage

import (
	"encoding/json"
	"git.ixarea.com/p2pNG/p2pNG-core/components/database"
	"git.ixarea.com/p2pNG/p2pNG-core/components/file_store"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func addLocalFile(c echo.Context) error {
	p := c.QueryParam("path")
	f, err := file_store.StatLocalFile(p, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 1, Message: "stat file error", Internal: err.Error()},
		)
	}
	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}
	fJson, err := json.Marshal(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 3, Message: "encoding json data error", Internal: err.Error()})
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("file"))
		return bucket.Put(f.Hash, fJson)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 4, Message: "write to database error", Internal: err.Error()})
	}
	return c.JSONBlob(http.StatusOK, fJson)
}
