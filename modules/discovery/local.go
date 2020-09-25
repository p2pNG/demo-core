package discovery

import (
	"github.com/micro/mdns"
	"os"
)

func LocalBroadcast(port int) (*mdns.Server, error) {
	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, _ := mdns.NewMDNSService(host, "_p2pNG._http", "", "", port, nil, info)
	return mdns.NewServer(&mdns.Config{Zone: service})
}

func LocalScan() (rt []mdns.ServiceEntry, err error) {
	entriesCh := make(chan *mdns.ServiceEntry, 64)
	err = mdns.Lookup("_p2pNG._http", entriesCh)
	if err != nil {
		return
	}
	//todo: Optimize
	for x := range entriesCh {
		rt = append(rt, *x)
	}
	close(entriesCh)
	return
}
