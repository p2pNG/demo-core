package main

import (
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/status"
	"github.com/davecgh/go-spew/spew"
	"net"
)

var PreDefined = net.TCPAddr{IP: net.ParseIP("192.168.133.1"), Port: 8444}

func main() {

	spew.Dump(status.GetNodeInfo(PreDefined))
	spew.Dump(discovery.ListAvailablePeers(PreDefined))

}
