package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mct-joken/jkojs-agent/docker"
	"github.com/mct-joken/jkojs-agent/types"
	"io"
	"net/http"
	"time"
)

func GetNewTask() (*Task, error) {
	resp, err := http.Get("http://localhost:3060/api/v2/submissions/tasks")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("no task to execute")
	}

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := Task{}
	err = json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func StartExec(task Task) error {
	cases := make([]types.ExecCases, len(task.Cases))
	for i, v := range task.Cases {
		cases[i] = types.ExecCases{
			Name: v.Name,
			File: []byte(v.Data),
		}
	}

	req := types.StartExecRequest{
		SubmissionID: task.ID,
		ProblemID:    task.ProblemID,
		Lang:         task.Lang,
		Code:         task.Code,
		Cases:        cases,
		Config: types.ExecConfig{
			TimeLimit:   task.Config.TimeLimit,
			MemoryLimit: task.Config.MemoryLimit,
		},
	}
	// SubmissionIDだけ空になるのでここで代入しておく
	res := &types.StartExecResponse{SubmissionID: task.ID}
	docker.Exec(req, res)
	fmt.Printf("%+#v\n", res)
	return nil
}

func AutoFetcher() {
	tick := time.NewTicker(500 * time.Millisecond)
	fmt.Println("start fetching...")
	lock := false
	for {
		select {
		case <-tick.C:
			if lock {
				continue
			}
			func() {
				lock = true
				defer func() {
					lock = false
				}()
				task, err := GetNewTask()
				if err != nil {
					fmt.Println(err)
				}
				err = StartExec(*task)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
	}
}

type Task struct {
	ID        string  `json:"ID"`
	ProblemID string  `json:"problemID"`
	Lang      string  `json:"lang"`
	Code      string  `json:"Code"`
	Cases     []Cases `json:"cases"`
	Config    Config  `json:"config"`
}

type Cases struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Config struct {
	TimeLimit   int `json:"timeLimit"`
	MemoryLimit int `json:"memoryLimit"`
}
