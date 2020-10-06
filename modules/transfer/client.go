package transfer

import (
	"encoding/base64"
	"encoding/json"
	"git.ixarea.com/p2pNG/p2pNG-core/components/file_store"
	"git.ixarea.com/p2pNG/p2pNG-core/components/request"
	"io/ioutil"
	"net"
	"net/url"
)

func GetSeed(tcpAddr net.TCPAddr, hash []byte) (seeds *file_store.FileInfo, err error) {
	endpoint := url.URL{Scheme: "https", Host: tcpAddr.String(), Path: "/transfer/seed/" + base64.RawURLEncoding.EncodeToString(hash)}

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
	seeds = new(file_store.FileInfo)
	err = json.Unmarshal(text, seeds)
	return
}
