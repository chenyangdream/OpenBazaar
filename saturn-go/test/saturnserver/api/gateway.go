package api

import (
	"github.com/op/go-logging"
	"github.com/Saturn/saturn-go/core"
	"net"
	"net/http"
)

var log = logging.MustGetLogger("api")

type Gateway struct {
	listener    net.Listener
	handler     http.Handler
	shutdownCh  chan struct{}
}

func NewGateway(n *core.SaturnNode, l net.Listener) (*Gateway, error){
	topMux := http.NewServeMux()

	jsonAPI, err := newJsonAPIHandler(n)

	if err != nil {
		return nil, err
	}

	topMux.Handle("/saturn/", jsonAPI)

	return &Gateway{
		listener:    l,
		handler:     topMux,
		shutdownCh:  make(chan struct{}),
	}, err

}

func (g *Gateway) Serve() error {
	return http.Serve(g.listener, g.handler)
}