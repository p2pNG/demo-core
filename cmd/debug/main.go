package main

import (
	"git.ixarea.com/p2pNG/p2pNG-core/modules/status"
	"github.com/davecgh/go-spew/spew"
	"net"
)

func main() {
	spew.Dump(status.GetNodeInfo(net.TCPAddr{IP: net.ParseIP("192.168.31.10"), Port: 8444}))
	spew.Dump(status.GetNodeInfo(net.TCPAddr{IP: net.ParseIP("192.168.31.10"), Port: 8443}))

}
