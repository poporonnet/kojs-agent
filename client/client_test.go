package client

import (
	"github.com/mct-joken/jkojs-agent/lib"
	"testing"
)

func TestStartExec(t *testing.T) {
	lib.InitLogger()
	lib.Config.ID = "c62af483410a6b204a12dba5de7e1d2b13dd400b7d756ab0f507b6fb2a64c6b5"
	task := Task{
		ID:        "200",
		ProblemID: "110",
		Lang:      "Python",
		Code:      "cHJpbnQoImhlbGxvIHdvcmxkXG4iKQ==",
		Cases: []Cases{
			{
				Name: "70.txt",
				Data: "hello\n",
			},
			{
				Name: "80.txt",
				Data: "1 2\n",
			},
			{
				Name: "90.txt",
				Data: "abc\n",
			},
			{
				Name: "100.txt",
				Data: "abc\nabp abc\n",
			},
		},
		Config: Config{
			TimeLimit:   512000,
			MemoryLimit: 512000,
		},
	}
	err := StartExec(task)
	if err != nil {
		return
	}
}
