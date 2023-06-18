package docker

import (
	"archive/tar"
	"encoding/base64"
	"github.com/mct-joken/jkojs-agent/types"
	"gotest.tools/v3/assert"
	"io"
	"testing"
)

func Test_decodeSourceCode(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("Hello world"))
	d := &types.StartExecRequest{
		SubmissionID: "0",
		ProblemID:    "0",
		Lang:         "Go",
		Code:         encoded,
		Cases:        nil,
		Config:       types.ExecConfig{},
	}
	_ = decodeSourceCode(d)
	assert.Equal(t, "Hello world", d.Code)
}

func Test_packSourceAndCases(t *testing.T) {
	cfg := types.TarFileDirectoryConfig{
		Payload: []types.TarFilePayload{
			{
				File: []byte("hello world"),
				Path: "./test/main.txt",
			},
			{
				File: []byte("echo $PATH"),
				Path: "./case/0001.txt",
			},
		},
	}
	res, _ := packSourceAndCases(cfg)
	reader := tar.NewReader(&res)

	i := 0
	for {
		file, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		assert.Equal(t, cfg.Payload[i].Path, file.Name)
		i++
	}
}
