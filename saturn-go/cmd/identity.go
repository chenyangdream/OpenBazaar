package ipfscmd

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/ipfs/go-ipfs/repo/config"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

func PeerIdFromPubKey(pk []byte) (string, error) {
	pubkey, err := libp2p.UnmarshalRsaPublicKey(pk)
	if err != nil {
		return "", err
	}
	id, err := peer.IDFromPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	return id.Pretty(), nil
}

func IdentityFromKey(privkey []byte) (config.Identity, error) {
	ident := config.Identity{}
	sk, err := libp2p.UnmarshalPrivateKey(privkey)
	if err != nil {
		return ident, err
	}
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(sk.GetPublic())
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}

func IdentityKeyFromSeed(seed []byte, bits int) ([]byte, error) {
	//hmac := hmac.New(sha256.New, []byte("Saturn seed"))
	//hmac.Write(seed)
	//reader := bytes.NewReader(hmac.Sum(nil))
	sk, _, err := libp2p.GenerateKeyPairWithReader(libp2p.RSA, bits, rand.Reader)
	if err != nil {
		return nil, err
	}
	encodedKey, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	return encodedKey, nil
}
