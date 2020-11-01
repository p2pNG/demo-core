package main

import (
	"github.com/p2pNG/core/modules/discovery"
	"net"
)

var PreDefined = net.TCPAddr{IP: net.ParseIP("192.168.133.1"), Port: 8444}

func main() {
	discovery.EnsureClientAlive()
}
