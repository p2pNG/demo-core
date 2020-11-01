package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/p2pNG/core/components/certificate"
	"github.com/p2pNG/core/components/file_store"
	"github.com/p2pNG/core/modules/discovery"
	"github.com/p2pNG/core/modules/status"
	"github.com/p2pNG/core/modules/transfer"
	"github.com/p2pNG/core/utils"
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
	_, _ = certificate.GetCertBundle("client", utils.GetHostname()+" Client")

	for cIdx := range clients {
		c := clients[cIdx]
		tcpAddr := net.TCPAddr{IP: c.Addr, Port: c.Port}
		utils.Log().Info("Found Local Service:", zap.String("addr", tcpAddr.String()))
		spew.Dump(status.GetNodeInfo(tcpAddr))
		//spew.Dump(manage.AddLocalFile(tcpAddr, "D:/temp/bank-proj/Release/package"))
		seeds, err := status.ListAvailableSeeds(tcpAddr)
		if err != nil {
			utils.Log().Error("failed to list peer owning seeds", zap.Error(err))
			continue
		}
		spew.Dump(seeds)
		for seedIdx := range seeds {
			seed, err := transfer.GetSeed(tcpAddr, seeds[seedIdx])
			if err != nil {
				utils.Log().Error("failed to get seed content", zap.Error(err))
				continue
			}
			spew.Dump(seed)
			for blockIdx := range seed.BlockHash {
				block, err := transfer.DownloadFileBlock(tcpAddr, seeds[seedIdx], seed.BlockHash[blockIdx])
				if err != nil {
					utils.Log().Error("failed to download file content", zap.Error(err))
				}
				spew.Dump(len(block) == file_store.DefaultHashBufferSize)
				break

			}
		}
	}

}
