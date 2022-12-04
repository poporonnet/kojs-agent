package types

type StartExecRequest struct {
	SubmissionID string      `json:"submissionID"`
	ProblemID    string      `json:"problemID"`
	Lang         string      `json:"lang"`
	Code         string      `json:"code"`
	Cases        []ExecCases `json:"cases"`
	Config       ExecConfig  `json:"config"`
}

type ExecCases struct {
	Name string `json:"name"`
	File []byte `json:"file"`
}

type ExecConfig struct {
	TimeLimit   int `json:"timeLimit"`
	MemoryLimit int `json:"memoryLimit"`
}

// コードとテストケースをWorkerに送るときの構造体
type TarFileDirectoryConfig struct {
	Payload []TarFilePayload
}

// Tarに入れるファイルの構造体
type TarFilePayload struct {
	File []byte
	Path string
}

// ここより下の2つはWorkerと共通
// ExecuteStatus 提出ごとのステータス
type StartExecResponse struct {
	SubmissionID string `json:"submissionID"` // 提出ID
	ProblemID    string `json:"problemID"`    // 問題ID
	LanguageType string `json:"languageType"` // 言語/処理系

	CompilerMessage     string `json:"compilerMessage"`     // コンパイラが出力した警告など
	CompileErrorMessage string `json:"compileErrorMessage"` // コンパイルエラー

	Results []CaseResult `json:"results"` // ケースごとのステータス
}

// CaseResult ケースごとのステータス
type CaseResult struct {
	Output      string `json:"output"`     // プログラム出力
	ExitStatus  int    `json:"exitStatus"` // 終了コード
	Duration    int    `json:"duration"`   // 実行時間
	MemoryUsage int    `json:"usage"`      // メモリ使用量
}

// 言語名からファイルを置く場所と名前を決定する
var LANGUAGE = map[string]string{
	"GCC":     "./test/main.c",
	"Clang":   "./test/main.c",
	"G++":     "./test/main.cpp",
	"Clang++": "./test/main.cpp",
	"Ruby":    "./built/main.rb", // ToDo: 決め打ちやめる
}

// ProblemConfig 問題ごとの設定
type ProblemConfig struct {
	ID          string   `json:"id"`
	TimeLimit   int      `json:"timeLimit"`   // 実行時間制限
	MemoryLimit int      `json:"memoryLimit"` // メモリ制限
	CaseFiles   []string `json:"caseFiles"`   // ケースファイルのファイルパス
}
