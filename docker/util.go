package docker

import (
	"encoding/base64"
	"fmt"

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

// Dockerから取ってきたファイルは先頭と末尾に要らない部分が大量にあるので,取り除く
func trimer(from []byte) []byte {
	// 先頭512byteはファイルヘッダー?
	tmp := from[512:]

	// 0x00が最初に現れたところで配列を切る
	cutPoint := 0
	for i, v := range tmp {
		if v == 0x00 {
			cutPoint = i
			break
		}
	}
	tmp = tmp[:cutPoint]
	fmt.Println(tmp)
	return tmp
}
