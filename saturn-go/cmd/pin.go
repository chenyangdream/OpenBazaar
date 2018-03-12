package ipfscmd

import "github.com/ipfs/go-ipfs/commands"

/* Recursively pin a directory given its hash. */
func Pin(ctx commands.Context, rootHash string) error {
	args := []string{"pin", "add", "/ipfs/" + rootHash}
	req, cmd, err := NewRequest(ctx, args)
	if err != nil {
		return err
	}
	res := commands.NewResponse(req)
	cmd.Run(req, res)
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

/* Recursively un-pin a directory given its hash.
   This will allow it to be garbage collected. */
func UnPin(ctx commands.Context, rootHash string) error {
	args := []string{"pin", "rm", "/ipfs/" + rootHash}
	req, cmd, err := NewRequest(ctx, args)
	if err != nil {
		return err
	}
	res := commands.NewResponse(req)
	cmd.Run(req, res)
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}