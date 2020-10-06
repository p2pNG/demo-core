package manage

import (
	"git.ixarea.com/p2pNG/p2pNG-core/components/request"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net"
	"net/url"
)

func AddLocalFile(tcpAddr net.TCPAddr, filepath string) (data string, err error) {
	qry := url.Values{"path": []string{filepath}}
	endpoint := url.URL{Scheme: "https", Host: tcpAddr.String(), Path: "/manage/add-file", RawQuery: qry.Encode()}
	spew.Dump(endpoint.String())

	client, err := request.GetDefaultHttpClient()
	if err != nil {
		return
	}

	resp, err := client.Get(endpoint.String())
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
