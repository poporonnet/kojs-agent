package manager

import (
	"context"
	"github.com/mct-joken/jkojs-agent/pkg/types"
)

type WorkerManager interface {
	// Start ワーカーを開始し、結果を受け取る
	Start(ctx context.Context, args StartWorkerArgs) (WorkerResponse, error)
}

// WorkerResponse ワーカーからの返答(実行の結果)
type WorkerResponse struct {
	SubmissionID string `json:"submissionID"`
	ProblemID    string `json:"problemID"`
	LanguageType string `json:"languageType"`

	CompilerMessage     string `json:"compilerMessage"`
	CompileErrorMessage string `json:"compileErrorMessage"`

	Results []CaseResult `json:"results"`
}

// CaseResult ケースごとのステータス
type CaseResult struct {
	CaseID      string `json:"caseID"`
	Output      string `json:"output"`
	ExitStatus  int    `json:"exitStatus"`
	Duration    int    `json:"duration"`
	MemoryUsage int    `json:"usage"`
}

// StartWorkerArgs ワーカー実行に必要な情報
type StartWorkerArgs struct {
	SubmissionID string         `json:"submissionID"`
	ProblemID    string         `json:"problemID"`
	Lang         types.LangCode `json:"lang"`
	Code         string         `json:"code"`
	Cases        []ExecCases    `json:"cases"`
	Config       ExecConfig     `json:"config"`
}

type ExecCases struct {
	Name string `json:"name"`
	File []byte `json:"file"`
}

type ExecConfig struct {
	TimeLimit   int `json:"timeLimit"`
	MemoryLimit int `json:"memoryLimit"`
}
