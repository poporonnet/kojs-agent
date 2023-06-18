package cmd

import (
	"fmt"
	dockerClient "github.com/docker/docker/client"
	"github.com/mct-joken/jkojs-agent/pkg/client"
	lib2 "github.com/mct-joken/jkojs-agent/pkg/lib"
	"github.com/mct-joken/jkojs-agent/pkg/manager/docker"
)

func Start(imageID, code, langID string) {
	lib2.InitLogger()
	lib2.Config.ID = imageID
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
	nclient, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	mng := docker.NewWorkerManager(nclient)

	err = client.StartExec(task, mng)
	if err != nil {
		fmt.Println(err)
		return
	}
}
