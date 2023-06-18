package cmd

import (
	"fmt"
	"github.com/mct-joken/jkojs-agent/client"
	"github.com/mct-joken/jkojs-agent/lib"
)

func Start(imageID, code, langID string) {
	lib.InitLogger()
	lib.Config.ID = imageID
	task := client.Task{
		ID:        "200",
		ProblemID: "110",
		Lang:      langID,
		Code:      code,
		Cases: []client.Cases{
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
		Config: client.Config{
			TimeLimit:   512000,
			MemoryLimit: 512000,
		},
	}
	err := client.StartExec(task)
	if err != nil {
		fmt.Println(err)
		return
	}
}
