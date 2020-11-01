package debug

import (
	"github.com/p2pNG/core/components/request"
	"io/ioutil"
	"net"
)

func DumpClientCertificate(tcpAddr net.TCPAddr) (data string, err error) {
	endpoint := "/debug/client-cert"

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

	data = string(text)
	return
}
