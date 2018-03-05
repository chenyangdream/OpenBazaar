package core

import (
	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("core")

var Node *SaturnNode

type SaturnNode struct {
	// Context for issuing IPFS commands
	Context commands.Context

	// IPFS node object
	IpfsNode *core.IpfsNode
}