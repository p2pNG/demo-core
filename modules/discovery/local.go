package discovery

import (
	"github.com/micro/mdns"
	"os"
)

// Attention: Should not Listen in the same goroutine with http service
func LocalBroadcast(port int) (*mdns.Server, error) {
	host, _ := os.Hostname()
	info := []string{"p2pNG Server Core"}
	service, _ := mdns.NewMDNSService(host, "_p2pNG._https", "", "", port, nil, info)
	return mdns.NewServer(&mdns.Config{Zone: service})
}

func LocalScan() ([]mdns.ServiceEntry, error) {
	var rt []mdns.ServiceEntry

	//todo: Optimize
	entriesCh := make(chan *mdns.ServiceEntry, 64)
	go func() {
		for x := range entriesCh {
			rt = append(rt, *x)
		}
	}()
	err := mdns.Lookup("_p2pNG._https", entriesCh)
	close(entriesCh)
	return rt, err
}
