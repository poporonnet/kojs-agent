package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Code-Hex/dd"
	"github.com/mct-joken/jkojs-agent/docker"
	"github.com/mct-joken/jkojs-agent/types"
	"io"
	"net/http"
	"time"
)

func GetNewTask() (Task, error) {
	resp, err := http.Get("https://ojs.laminne33569.net/api/v2/submissions/tasks")
	if err != nil {
		return Task{}, err
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return Task{}, errors.New("no task to execute")
	}

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return Task{}, err
	}
	data := Task{}
	err = json.Unmarshal(d, &data)
	if err != nil {
		return Task{}, err
	}
	fmt.Println(dd.Dump(data))
	return data, nil
}

func NotifyTaskFinished(t types.StartExecResponse) error {
	results := make([]CreateSubmissionResults, len(t.Results))
	for i, v := range t.Results {
		results[i] = CreateSubmissionResults{
			CaseName:   v.CaseID,
			Output:     v.Output,
			ExitStatus: v.ExitStatus,
			Duration:   v.Duration,
			Usage:      v.MemoryUsage,
		}
	}
	b := Request{
		SubmissionID:        t.SubmissionID,
		ProblemID:           t.ProblemID,
		LanguageType:        t.LanguageType,
		CompilerMessage:     t.CompilerMessage,
		CompileErrorMessage: t.CompileErrorMessage,
		Results:             results,
	}
	req, _ := json.Marshal(b)
	a := bytes.NewBuffer(req)
	_, err := http.Post("https://ojs.laminne33569.net/api/v2/submissions/tasks", "application/json", a)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
	res.SubmissionID = task.ID
	fmt.Println(*res)
	return NotifyTaskFinished(*res)
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
					return
				}
				err = StartExec(task)
				if err != nil {
					fmt.Println(err)
					return
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

type Request struct {
	SubmissionID        string                    `json:"submissionID"`
	ProblemID           string                    `json:"problemID"`
	LanguageType        string                    `json:"languageType"`
	CompilerMessage     string                    `json:"compilerMessage"`
	CompileErrorMessage string                    `json:"compileErrorMessage"`
	Results             []CreateSubmissionResults `json:"results"`
}

type CreateSubmissionResults struct {
	CaseName   string `json:"caseName"`
	Output     string `json:"output"`
	ExitStatus int    `json:"exitStatus"`
	Duration   int    `json:"duration"`
	Usage      int    `json:"usage"`
}
