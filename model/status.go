package model

type NodeInfo struct {
	Name      string
	Version   string
	BuildName string
}

type RegReqNodeInfo struct {
	NodeInfo
	Port int
}

type RegNodeInfo struct {
	NodeInfo
	Port    int
	Address string
}
