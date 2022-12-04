package main

import (
	"github.com/mct-joken/jkojs-agent/docker"
	"github.com/mct-joken/jkojs-agent/lib"
	"github.com/mct-joken/jkojs-agent/types"
)

func main() {
	lib.InitLogger()
	a := &types.StartExecResponse{}
	docker.Exec(types.StartExecRequest{
		SubmissionID: "",
		ProblemID:    "112233",
		Lang:         "Clang++",
		Code:         "I2luY2x1ZGUgPGlvc3RyZWFtPgoKdXNpbmcgbmFtZXNwYWNlIHN0ZDsKCmludCBtYWluKCkgewogICAgY291dCA8PCAiSGVsbG8gV29ybGQgQysrIiA8PCBlbmRsOwogICAgcmV0dXJuIDA7Cn0K",
		Cases: []types.ExecCases{
			{
				Name: "test.txt",
				File: []byte("Hello World C++"),
			},
		},
		Config: types.ExecConfig{
			TimeLimit:   2000,
			MemoryLimit: 512,
		},
	}, a)
	lib.Logger.Sugar().Debugln(a)
	// api.StartServer()
}
