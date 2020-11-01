package transfer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/p2pNG/core/components/database"
	"github.com/p2pNG/core/model"
	bolt "go.etcd.io/bbolt"
	"io"
	"net/http"
	"os"
)

func getSeed(c echo.Context) error {
	hash, err := base64.RawURLEncoding.DecodeString(c.Param("hash"))
	if err != nil || len(hash) != 64 {
		if err == nil {
			err = errors.New("incorrect hash length")
		}
		return c.JSON(http.StatusBadRequest,
			model.StandardError{Code: 6, Message: "parse seed hash error", Internal: err.Error()})
	}

	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}

	var rawData []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("file"))
		rawData = bucket.Get(hash)
		return nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 5, Message: "read from database error", Internal: err.Error()})
	}
	if rawData == nil {
		return c.JSON(http.StatusNotFound,
			model.StandardError{Code: 8, Message: "no such seed"})
	}
	var data model.LocalFileInfo
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 9, Message: "decode data error", Internal: err.Error()})
	}
	return c.JSON(http.StatusOK, data.FileInfo)

}

func downloadFileBlock(c echo.Context) error {
	hash, err := base64.RawURLEncoding.DecodeString(c.Param("hash"))
	if err != nil || len(hash) != 64 {
		if err == nil {
			err = errors.New("incorrect hash length")
		}
		return c.JSON(http.StatusBadRequest,
			model.StandardError{Code: 6, Message: "parse seed hash error", Internal: err.Error()})
	}
	block, err := base64.RawURLEncoding.DecodeString(c.Param("block"))
	if err != nil || len(block) != 32 {
		if err == nil {
			err = errors.New("incorrect hash length")
		}
		return c.JSON(http.StatusBadRequest,
			model.StandardError{Code: 13, Message: "parse block hash error", Internal: err.Error()})
	}

	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}

	var rawData []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("file"))
		rawData = bucket.Get(hash)
		return nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 5, Message: "read from database error", Internal: err.Error()})
	}
	if rawData == nil {
		return c.JSON(http.StatusNotFound,
			model.StandardError{Code: 8, Message: "no such seed"})
	}
	var lf model.LocalFileInfo
	err = json.Unmarshal(rawData, &lf)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 9, Message: "decode data error", Internal: err.Error()})
	}
	blockIdx := -1
	for _blockIdx := range lf.BlockHash {
		if bytes.Equal(lf.BlockHash[_blockIdx], block) {
			blockIdx = _blockIdx
			break
		}
	}
	if blockIdx == -1 {
		return c.JSON(http.StatusNotFound,
			model.StandardError{Code: 10, Message: "no such block in this seed"})
	}
	startPos := int64(blockIdx) * lf.BlockSize
	endPos := int64(blockIdx+1) * lf.BlockSize
	if endPos > lf.Size {
		endPos = lf.Size
	}
	f, err := openFile(lf.Path, startPos)
	if blockIdx == -1 {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 11, Message: "open file failed"})
	}

	resp := c.Response()
	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.WriteHeader(http.StatusOK)
	_, err = io.CopyN(resp, f, endPos-startPos)
	_ = f.Close()
	return err
}
func openFile(filepath string, startPos int64) (f *os.File, err error) {
	f, err = os.Open(filepath)
	if err != nil {
		return
	}
	_, err = f.Seek(startPos, io.SeekStart)
	if err != nil {
		return
	}
	return f, err
}
