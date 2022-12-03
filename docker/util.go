package docker

import (
	"encoding/base64"

	"github.com/mct-joken/jkojs-agent/types"
)

func decodeSourceCode(arg *types.StartExecRequest) error {
	enc, err := base64.StdEncoding.DecodeString(arg.Code)
	if err != nil {
		return err
	}
	arg.Code = string(enc)
	return nil
}
