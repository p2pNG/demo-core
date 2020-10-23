package manage

import (
	"encoding/json"
	"git.ixarea.com/p2pNG/p2pNG-core/components/database"
	"git.ixarea.com/p2pNG/p2pNG-core/components/file_store"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
	"net"
	"net/http"
)

func addLocalFile(c echo.Context) error {
	if !net.ParseIP(c.RealIP()).IsLoopback() {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 16, Message: "should called by loopback ip"},
		)
	}
	p := c.QueryParam("path")
	f, err := file_store.StatLocalFile(p, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 1, Message: "stat file error", Internal: err.Error()},
		)
	}
	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}
	fJson, err := json.Marshal(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 3, Message: "encoding json data error", Internal: err.Error()})
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("file"))
		return bucket.Put(f.Hash, fJson)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 4, Message: "write to database error", Internal: err.Error()})
	}
	return c.JSONBlob(http.StatusOK, fJson)
}
