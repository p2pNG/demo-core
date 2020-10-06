package transfer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"git.ixarea.com/p2pNG/p2pNG-core/components/database"
	"git.ixarea.com/p2pNG/p2pNG-core/components/file_store"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func getSeed(c echo.Context) error {
	hash, err := base64.RawURLEncoding.DecodeString(c.Param("hash"))
	if err != nil || len(hash) != 64 {
		if err == nil {
			err = errors.New("incorrect hash length")
		}
		return c.JSON(http.StatusBadRequest,
			utils.StandardError{Code: 6, Message: "parse seed hash error", Internal: err.Error()})
	}

	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}
	var rawData []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("file"))
		rawData = bucket.Get(hash)
		return nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 5, Message: "read from database error", Internal: err.Error()})
	}
	if rawData == nil {
		return c.JSON(http.StatusNotFound,
			utils.StandardError{Code: 8, Message: "no such seed"})
	}
	var data file_store.LocalFileInfo
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.StandardError{Code: 9, Message: "decode data error", Internal: err.Error()})
	}
	return c.JSON(http.StatusOK, data.FileInfo)

}
