package discovery

import (
	"bytes"
	"encoding/json"
	core "git.ixarea.com/p2pNG/p2pNG-core"
	"git.ixarea.com/p2pNG/p2pNG-core/components/request"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"io/ioutil"
	"net"
	"net/http"
)

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
