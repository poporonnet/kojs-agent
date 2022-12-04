package docker

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mct-joken/jkojs-agent/lib"
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

// 構造体の内容をtarにまとめる
func packSourceAndCases(conf types.TarFileDirectoryConfig) (bytes.Buffer, error) {
	// ローカルにファイルを作りたくないのでメモリ上に持つ
	var tarFile bytes.Buffer
	writer := tar.NewWriter(&tarFile)
	defer writer.Close()

	for _, v := range conf.Payload {
		if err := writer.WriteHeader(&tar.Header{
			Name:    v.Path,
			Mode:    int64(0666), // ToDo: パーミッションの検討
			ModTime: time.Now(),
			Size:    int64(len(v.File)),
		}); err != nil {
			return tarFile, err
		}
		writer.Write(v.File)
	}

	// f, _ := os.OpenFile("test.tar", os.O_RDWR|os.O_CREATE, 0666)
	// a, _ := io.ReadAll(&tarFile)
	// f.Write(a)
	// f.Close()

	return tarFile, nil
}

// tarにまとめるファイルを構造体に書き込む
func preparePacking(req types.StartExecRequest) types.TarFileDirectoryConfig {
	config := types.TarFileDirectoryConfig{}

	// ケースファイルを詰める
	for _, v := range req.Cases {
		payload := types.TarFilePayload{
			File: v.File,
			Path: "./case/" + v.Name,
		}
		config.Payload = append(config.Payload, payload)
	}

	// ソースコードを詰める
	programFile := types.TarFilePayload{
		Path: types.LANGUAGE[req.Lang],
		File: []byte(req.Code),
	}
	lib.Logger.Sugar().Debugf("%v", programFile)
	config.Payload = append(config.Payload, programFile)

	// 問題ごとの実行時設定を取る
	c := types.ProblemConfig{
		ID:          req.ProblemID,
		TimeLimit:   req.Config.TimeLimit,
		MemoryLimit: req.Config.MemoryLimit,
	}
	for _, v := range req.Cases {
		c.CaseFiles = append(c.CaseFiles, v.Name)
	}
	execFileCfg, _ := json.Marshal(c)

	execConfigFile := types.TarFilePayload{
		File: execFileCfg,
		Path: "./case/" + req.ProblemID + ".json",
	}
	config.Payload = append(config.Payload, execConfigFile)

	return config
}
