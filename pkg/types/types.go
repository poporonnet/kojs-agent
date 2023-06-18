package types

import "errors"

type Config struct {
	ID string `yaml:"imageID"`
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

// CaseResult ケースごとのステータス
type CaseResult struct {
	CaseID      string `json:"caseID"`     // ケースID
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
	"Ruby":    "./test/main.rb",
	"Go":      "./test/main.go",
	"Python3": "./test/main.py",
}

type LangCode struct {
	code int
}

func (c LangCode) Code() int {
	return c.code
}

func (c LangCode) ToString() string {
	switch c.code {
	case GCC:
		return "GCC"
	case Clang:
		return "Clang"
	case GXX:
		return "G++"
	case ClangXX:
		return "Clang++"
	case Ruby:
		return "Ruby"
	case Go:
		return "Go"
	case Python3:
		return "Python3"
	default:
		return ""
	}
}

func NewLangCode(code string) (LangCode, error) {
	switch code {
	case "GCC":
		return LangCode{
			code: GCC,
		}, nil
	case "Clang":
		return LangCode{
			code: Clang,
		}, nil
	case "G++":
		return LangCode{
			code: GXX,
		}, nil
	case "Clang++":
		return LangCode{
			code: ClangXX,
		}, nil
	case "Ruby":
		return LangCode{
			code: Ruby,
		}, nil
	case "Go":
		return LangCode{
			code: Go,
		}, nil
	case "Python3":
		return LangCode{
			code: Python3,
		}, nil
	default:
		return LangCode{}, errors.New("no such lang")
	}
}

const (
	GCC = iota
	Clang
	GXX
	ClangXX
	Ruby
	Go
	Python3
)

// ProblemConfig 問題ごとの設定
type ProblemConfig struct {
	ID          string   `json:"id"`
	TimeLimit   int      `json:"timeLimit"`   // 実行時間制限
	MemoryLimit int      `json:"memoryLimit"` // メモリ制限
	CaseFiles   []string `json:"caseFiles"`   // ケースファイルのファイルパス
}
