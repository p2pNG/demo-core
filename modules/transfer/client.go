package transfer

import (
	"encoding/base64"
	"encoding/json"
	"git.ixarea.com/p2pNG/p2pNG-core/components/request"
	"git.ixarea.com/p2pNG/p2pNG-core/model"
	"io/ioutil"
	"net"
	"net/url"
)

func GetSeed(tcpAddr net.TCPAddr, hash []byte) (seeds *model.FileInfo, err error) {
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
	seeds = new(model.FileInfo)
	err = json.Unmarshal(text, seeds)
	return
}

func DownloadFileBlock(tcpAddr net.TCPAddr, hash []byte, block []byte) (data []byte, err error) {
	endpoint := url.URL{Scheme: "https", Host: tcpAddr.String(),
		Path: "/transfer/seed/" + base64.RawURLEncoding.EncodeToString(hash) +
			"/" + base64.RawURLEncoding.EncodeToString(block),
	}

	client, err := request.GetDefaultHttpClient()
	if err != nil {
		return
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(resp.Body)

	return
}
