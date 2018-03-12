package ipfscmd

import (
	"errors"
	"github.com/ipfs/go-ipfs/commands"
)

var getErr = errors.New(`get file failed`)

/*Get file to outputPath by fileHash from ipfs network.*/
func Get(ctx commands.Context, fileHash , outputPath string) error {
	args := []string{"get", fileHash, "-o", outputPath}
	req, cmd, err := NewRequest(ctx, args)
	if err != nil {
		return err
	}
	res := commands.NewResponse(req)
	cmd.PreRun(req)
	cmd.Run(req, res)
	cmd.PostRun(req, res)
	if res.Error() != nil {
		return res.Error()
	}

	return nil
}