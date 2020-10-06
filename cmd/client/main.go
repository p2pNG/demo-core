package main

import (
	"git.ixarea.com/p2pNG/p2pNG-core/components/certificate"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/discovery"
	"git.ixarea.com/p2pNG/p2pNG-core/modules/status"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"net"
)

func main() {
	utils.Log().Info("Hello")
	clients, err := discovery.LocalScan()
	if err != nil {
		utils.Log().Error("scan local peer failed", zap.Error(err))
	}
	if err != nil {
		utils.Log().Fatal("open database", zap.Error(err))
	}

	_, err = certificate.GetCert("client", utils.GetHostname()+" Client")
	if err != nil {
		utils.Log().Fatal("generate local certificate failed", zap.Error(err))
	}

	for cIdx := range clients {
		c := clients[cIdx]
		tcpAddr := net.TCPAddr{IP: c.Addr, Port: c.Port}
		utils.Log().Info("Found Local Service:", zap.String("addr", tcpAddr.String()))
		spew.Dump(status.GetNodeInfo(tcpAddr))
		//spew.Dump(manage.AddLocalFile(tcpAddr, "D:/temp/bank-proj/Release/package"))
		spew.Dump(status.ListAvailableSeeds(tcpAddr))
	}

}
