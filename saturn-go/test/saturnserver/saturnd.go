package main

import (
	"fmt"
	"github.com/Saturn/saturn-go"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"github.com/Saturn/saturn-go/core"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	manet "gx/ipfs/QmX3U3YXCQ6UYBxq2LVWF8dARS1hPUTEYLrSx654Qyxyw6/go-multiaddr-net"
	"github.com/Saturn/saturn-go/test/saturnserver/api"
)

var log = logging.MustGetLogger("main")

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Noticef("Received %s\n", sig)
			log.Info("Saturn Server shutting down...")

			err := ipfs.Shutdown()
			if err != nil {
				log.Infof("%s", err.Error())
				os.Exit(1)
			}

			os.Exit(1)
		}
	}()

	err := ipfs.Startup(nil, "")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	node, err:= core.GetSaturnNode()
	if err != nil {
		fmt.Println("Saturn Node is nil Please core ifs.Startup create Saturn Node first!")
		os.Exit(1)
	}


	//fmt.Print("Press 'Enter' to continue ...")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	gateway, err := newHTTPGateway(node)

	if err != nil {
		fmt.Println("new HTTP Gateway error")
		os.Exit(1)
	}

	gateway.Serve()
}

// Collects options, creates listener, prints status message and starts serving requests
func newHTTPGateway(node *core.SaturnNode) (*api.Gateway, error) {

	address := "/ip4/0.0.0.0/tcp/4002"
	// Create a network listener
	gatewayMaddr, err := ma.NewMultiaddr(address)
	if err != nil {
		return nil, fmt.Errorf("newHTTPGateway: invalid gateway address: %q (err: %s)", address, err)
	}
	var gwLis manet.Listener
	gwLis, err = manet.Listen(gatewayMaddr)
	if err != nil {
		return nil, fmt.Errorf("newHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
	}

	// We might have listened to /tcp/0 - let's see what we are listing on
	gatewayMaddr = gwLis.Multiaddr()
	log.Infof("Gateway/API server listening on %s\n", gatewayMaddr)

	if err != nil {
		return nil, fmt.Errorf("newHTTPGateway: ConstructNode() failed: %s", err)
	}

	return api.NewGateway(node, gwLis.NetListener())
}