package ipfs

import (
	"context"
	"github.com/Saturn/saturn-go/core"
	"github.com/Saturn/saturn-go/cmd"
	"crypto/rsa"
)

type PeerInfo struct {
	PeerId     string
	PrivateKey *rsa.PrivateKey
}

/*Startup the saturn and ipfs library.*/
func Startup(bootstrapPeers []string, cachePath string) error {
	cctx, cancel := context.WithCancel(context.Background())
	err := core.Start(cctx, cachePath, cancel)
	if err != nil {
		return err
	}
	return nil
}

/*Shutdown the saturn node and ipfs node.*/
func Shutdown() error {
	node, err := core.GetSaturnNode()
	if err != nil {
		return err
	}
	node.CancelCtx()
	return nil
}

/*Returns the Saturn node PeerInfo*/
func GetPeerInfo() (*PeerInfo, error) {
	node, err := core.GetSaturnNode()
	if err != nil {
		return nil, err
	}

	peerid, err := node.GetPeerId()
	if err != nil {
		return nil, err
	}

	privkey, err := node.GetPriKey()
	if err != nil {
		return nil, err
	}

	peerInfo := &PeerInfo{
		PeerId:     peerid,
		PrivateKey: privkey,
	}

	return peerInfo, nil
}

/*Get peer id from rsa-publickey.*/
func PeerIdFromPubKey(pkbytes []byte) (string, error) {
	peerid, err := ipfscmd.PeerIdFromPubKey(pkbytes)
	if err != nil {
		return "", err
	}
	return peerid, nil
}

/*Add a file to ipfs. */
func AddFile(filePath string) (fileHash string, err error) {
	node, err := core.GetSaturnNode()
	if err != nil {
		return "", err
	}
	return ipfscmd.AddFile(node.Context, filePath)
}

/*Pin the specified file by fileHash.*/
func PinFile(fileHash string) error {
	node, err := core.GetSaturnNode()
	if err != nil {
		return err
	}
	return ipfscmd.Pin(node.Context, fileHash)
}

/*Unpin the specified file by fileHash.*/
func UnpinFile(fileHash string) error {
	node, err := core.GetSaturnNode()
	if err != nil {
		return err
	}
	return ipfscmd.UnPin(node.Context, fileHash)
}

/*Get(Download) the specified file by fileHash to outDir*/
func GetFile(fileHash string, outDir string, peers []string) error {
	node, err := core.GetSaturnNode()
	if err != nil {
		return err
	}
	return ipfscmd.Get(node.Context, fileHash, outDir)
}

/**
  Get the download progress.
  Parameters:
    fileHash: the file hash.
  Return:
    fileSize: the total size of the file.
    downloadedsize: the downloaded size.
 */
func GetDownloadProgress(fileHash string) (fileSize, downloadedSize uint64) {
	return 0, 0
}

/**
  Get the uploaded data size of the specified peer for a file.
  Parameters:
    fileHash: the file hash.
    peerId: the peer id.
  Return:
    uploaded data size (Bytes)
 */
func StatTraffic(fileHash string, peerId string) int64 {
	return 0
}