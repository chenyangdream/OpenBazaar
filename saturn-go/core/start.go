package core

import (
	"context"
	"fmt"
	"github.com/Saturn/saturn-go/repo"
	"github.com/ipfs/go-ipfs/commands"
	ipfscore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/namesys"
	namepb "github.com/ipfs/go-ipfs/namesys/pb"
	ipfsrepo "github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	lockfile "github.com/ipfs/go-ipfs/repo/fsrepo/lock"
	"github.com/ipfs/go-ipfs/thirdparty/ds-help"
	routing "gx/ipfs/QmPR2JzfKd9poHx9XBhzoFeBBC31ZM3W5iUPKJZWyaoZZm/go-libp2p-routing"
	dht "gx/ipfs/QmUCS9EnqNq1kCnJds2eLDypBiS21aSiCf1MVzSUVB9TGA/go-libp2p-kad-dht"
	dhtutil "gx/ipfs/QmUCS9EnqNq1kCnJds2eLDypBiS21aSiCf1MVzSUVB9TGA/go-libp2p-kad-dht/util"
	proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	p2phost "gx/ipfs/QmaSxYRuMq4pkpBBG2CYaRrPx2z7NmMVEs34b9g61biQA6/go-libp2p-host"
	recpb "gx/ipfs/QmbxkgUceEcuSZ4ZdBA3x74VUDSSYjHYmmeEqkjxbtZ6Jg/go-libp2p-record/pb"
	"time"
	"path/filepath"
	"os"
)

func Start(cctx context.Context, repoPath string, cancel context.CancelFunc) error {
	fmt.Println("Saturn Server v0.1.0.........")

	// Set repo path
	repoPath, err := repo.GetRepoPath()
	if err != nil {
		return err
	}

	repoLockFile := filepath.Join(repoPath, lockfile.LockFile)
	os.Remove(repoLockFile)

	isTestnet := false
	err = InitializeRepo(repoPath, "", "", isTestnet, time.Now())
	if err != nil && err != repo.ErrRepoExists {
		return err
	}

	err = CheckAndSetUlimit()
	if err != nil {
		return err
	}

	// IPFS node setup
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		log.Error(err)
		return err
	}


	cfg, err := r.Config()
	if err != nil {
		log.Error(err)
		return err
	}

	//config Swarm Address commented by ChenYang 2018-03-02 15:51:00
	cfg.Addresses.Swarm = append(cfg.Addresses.Swarm, "/ip4/0.0.0.0/tcp/4001")
	cfg.Addresses.Swarm = append(cfg.Addresses.Swarm, "/ip6/::/tcp/4001")

	// If we're only using Tor set the proxy dialer and dns resolver
	dnsResolver := namesys.NewDNSResolver()

	ncfg := &ipfscore.BuildCfg{
		Repo:   r,
		Online: true,
		ExtraOpts: map[string]bool{
			"mplex": true,
		},
		DNSResolver: dnsResolver,
		Routing:     DHTOption,
	}

	nd, err := ipfscore.NewNode(cctx, ncfg)
	if err != nil {
		log.Error(err)
		return err
	}

	ctx := commands.Context{}
	ctx.Online = true
	ctx.ConfigRoot = repoPath
	ctx.LoadConfig = func(path string) (*config.Config, error) {
		return fsrepo.ConfigAt(repoPath)
	}
	ctx.ConstructNode = func() (*ipfscore.IpfsNode, error) {
		return nd, nil
	}

	// Set IPNS query size
	querySize := cfg.Ipns.QuerySize
	if querySize <= 20 && querySize > 0 {
		dhtutil.QuerySize = int(querySize)
	} else {
		dhtutil.QuerySize = 16
	}
	namesys.UsePersistentCache = cfg.Ipns.UsePersistentCache

	log.Info("Peer ID: ", nd.Identity.Pretty())
	fmt.Println("Peer ID: ", nd.Identity.Pretty())

	// Get current directory root hash
	_, ipnskey := namesys.IpnsKeysForID(nd.Identity)
	ival, hasherr := nd.Repo.Datastore().Get(dshelp.NewKeyFromBinary([]byte(ipnskey)))
	if hasherr != nil {
		log.Error(hasherr)
		return err
	}
	val := ival.([]byte)
	dhtrec := new(recpb.Record)
	proto.Unmarshal(val, dhtrec)
	e := new(namepb.IpnsEntry)
	proto.Unmarshal(dhtrec.GetValue(), e)

	node = &SaturnNode{
		Context:             ctx,
		IpfsNode:            nd,
		CancelCtx:           cancel,
		RepoPath:            repoPath,
	}

	return err
}

var DHTOption ipfscore.RoutingOption = constructDHTRouting

func constructDHTRouting(ctx context.Context, host p2phost.Host, dstore ipfsrepo.Datastore) (routing.IpfsRouting, error) {
	dhtRouting := dht.NewDHT(ctx, host, dstore)
	dhtRouting.Validator[ipfscore.IpnsValidatorTag] = namesys.IpnsRecordValidator
	dhtRouting.Selector[ipfscore.IpnsValidatorTag] = namesys.IpnsSelectorFunc
	return dhtRouting, nil
}

func InitializeRepo(dataDir, password, mnemonic string, testnet bool, creationDate time.Time) error {

	// Initialize the IPFS repo if it does not already exist
	err := repo.DoInit(dataDir, 4096, testnet, password, mnemonic, creationDate)
	if err != nil {
		return err
	}
	return nil
}
