package docker

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mct-joken/jkojs-agent/pkg/lib"
	"github.com/mct-joken/jkojs-agent/pkg/manager"
	"github.com/mct-joken/jkojs-agent/pkg/types"
	"time"
)

func decodeSourceCode(arg *manager.StartWorkerArgs) error {
	enc, err := base64.StdEncoding.DecodeString(arg.Code)
	if err != nil {
		lib.Logger.Sugar().Error(err, arg.Code)
		return err
	}
	arg.Code = string(enc)
	return nil
}

// 構造体の内容をtarにまとめる
func packSourceAndCases(conf types.TarFileDirectoryConfig) (bytes.Buffer, error) {
	// ローカルにファイルを作りたくないのでメモリ上に持つ
	var tarFile bytes.Buffer
	writer := tar.NewWriter(&tarFile)
	defer func() {
		_ = writer.Close()
	}()

	for _, v := range conf.Payload {
		if err := writer.WriteHeader(&tar.Header{
			Name:    v.Path,
			Mode:    int64(0666), // ToDo: パーミッションの検討
			ModTime: time.Now(),
			Size:    int64(len(v.File)),
		}); err != nil {
			return tarFile, err
		}
		_, err := writer.Write(v.File)
		if err != nil {
			fmt.Println(err)
			return bytes.Buffer{}, err
		}
	}

	return tarFile, nil
}

// tarにまとめるファイルを構造体に書き込む
func preparePacking(req manager.StartWorkerArgs) types.TarFileDirectoryConfig {
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
		Path: types.LANGUAGE[req.Lang.ToString()],
		File: []byte(req.Code),
	}
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
