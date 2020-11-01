package status

import (
	"encoding/json"
	"github.com/p2pNG/core/components/request"
	"github.com/p2pNG/core/model"
	"io/ioutil"
	"net"
)

func GetNodeInfo(tcpAddr net.TCPAddr) (info *model.NodeInfo, err error) {
	endpoint := "/status/info"

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
	info = new(model.NodeInfo)
	err = json.Unmarshal(text, info)
	return
}

func ListAvailableSeeds(tcpAddr net.TCPAddr) (seeds [][]byte, err error) {
	endpoint := "/status/seeds"

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
	seeds = [][]byte{}
	err = json.Unmarshal(text, &seeds)
	return
}
