package discovery

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/p2pNG/core/components/database"
	"github.com/p2pNG/core/model"
	"github.com/p2pNG/core/modules/status"
	"github.com/p2pNG/core/utils"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"net"
	"net/http"
)

func registerClient(c echo.Context) error {
	req := new(model.RegReqNodeInfo)
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			model.StandardError{Code: 6, Message: "parse request parameter error", Internal: err.Error()})
	}
	peer := net.TCPAddr{IP: net.ParseIP(c.RealIP()), Port: req.Port}
	st, err := status.GetNodeInfo(peer)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable,
			model.StandardError{Code: 13, Message: "checking connection failed", Internal: err.Error()})
	}
	if *st != req.NodeInfo {
		return c.JSON(http.StatusNotAcceptable,
			model.StandardError{Code: 13, Message: "checking connection failed, not the same"})
	}
	data, err := json.Marshal(st)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 3, Message: "encoding json data error", Internal: err.Error()})
	}
	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}
	err = db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("discovery_registry")).Put([]byte(peer.String()), data)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 4, Message: "write to database error", Internal: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func listAvailablePeers(c echo.Context) error {
	db, err := database.GetDBEngine()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 2, Message: "connect to database error", Internal: err.Error()})
	}
	var data []string
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("discovery_registry"))
		return bucket.ForEach(func(k, _ []byte) error {
			data = append(data, string(k))
			return nil
		})
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			model.StandardError{Code: 5, Message: "read from database error", Internal: err.Error()})
	}
	return c.JSON(http.StatusOK, data)

}

func EnsureClientAlive() (err error) {
	db, err := database.GetDBEngine()
	if err != nil {
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("discovery_registry"))
		return bucket.ForEach(func(k, _ []byte) error {
			addr, err := net.ResolveTCPAddr("tcp", string(k))
			if err != nil {
				return err
			}
			_, err = status.GetNodeInfo(*addr)
			//todo: Reserve Offline Clients For
			if err != nil {
				utils.Log().Debug("Deleting", zap.Binary("key", k))
				_ = bucket.Delete(k)
			}
			return nil
		})
	})
	return
}
