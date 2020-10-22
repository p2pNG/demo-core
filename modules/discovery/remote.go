package discovery

import (
	"bytes"
	"encoding/json"
	core "git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/components/database"
	"git.ixarea.com/p2pNG/p2pNG-core/components/request"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/status"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
)

func GetRouter(e *echo.Group) {
	g := e.Group("/discovery")
	g.POST("/register", registerClient)
	g.GET("/peers", listAvailablePeers)
}

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

func RegisterClient(addr net.TCPAddr, listenPort int) (success bool, err error) {
	success = false
	message, err := json.Marshal(model.RegReqNodeInfo{NodeInfo: model.NodeInfo{
		Name:      utils.GetHostname(),
		Version:   core.GetVersionTag(),
		BuildName: core.GetBuildName(),
	}, Port: listenPort})

	endpoint := "/discovery/register"
	client, err := request.GetDefaultHttpClient()
	if err != nil {
		return
	}

	resp, err := client.Post("https://"+addr.String()+endpoint, "application/json", bytes.NewReader(message))
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	stdErr := new(model.StandardError)
	err = json.Unmarshal(data, stdErr)
	if err != nil {
		return
	}
	return false, stdErr
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

func ListAvailablePeers(tcpAddr net.TCPAddr) (seeds []string, err error) {
	endpoint := "/discovery/peers"

	client, err := request.GetDefaultHttpClient()
	if err != nil {
		return
	}

	resp, err := client.Get("https://" + tcpAddr.String() + endpoint)
	if err != nil {
		return
	}

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	seeds = []string{}
	err = json.Unmarshal(text, &seeds)
	return
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
