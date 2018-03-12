package core

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/op/go-logging"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)
var log = logging.MustGetLogger("core")

var ErrSaturnNodeNil = errors.New("Saturn node is nil. You should create Saturn node first!")
var ErrIpfsNodeNil = errors.New("Ipfs node is nil. You should create Ipfs node first!")

var node *SaturnNode

type SaturnNode struct {
	// Context for issuing IPFS commands
	Context commands.Context

	// IPFS node object
	IpfsNode *core.IpfsNode

	// when Saturn exits
	CancelCtx context.CancelFunc

	/* The roothash of the node directory inside the saturn repo.
	   This directory hash is published on IPNS at our peer ID making
	   the directory publicly viewable on the network. */
	RootHash string

	// The path to the saturn repo in the file system
	RepoPath string
}

func GetSaturnNode() (*SaturnNode, error) {
	if node == nil {
		return nil, ErrSaturnNodeNil
	}
	return node, nil
}

func (node *SaturnNode) GetContext() (context.Context, error) {
	if node.IpfsNode == nil {
		return nil, ErrIpfsNodeNil
	}
	return node.IpfsNode.Context(), nil
}

func (node *SaturnNode) GetPeerId() (string, error) {
	if node.IpfsNode == nil {
		return "", ErrIpfsNodeNil
	}
	return node.IpfsNode.Identity.Pretty(), nil
}

func (node *SaturnNode) GetPriKey() (*rsa.PrivateKey, error) {
	if node.IpfsNode == nil {
		return nil, ErrIpfsNodeNil
	}

	rsaprikey, ok := node.IpfsNode.PrivateKey.(*libp2p.RsaPrivateKey)
	if !ok {
		return nil, errors.New("IpfsNode get rsa PrivateKey failed")
	}

	return rsaprikey.GetPrivateKey(), nil
}