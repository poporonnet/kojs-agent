package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mct-joken/jkojs-agent/lib"
	jkojsTypes "github.com/mct-joken/jkojs-agent/types"
)

type cli struct {
	c         *client.Client
	container container.ContainerCreateCreatedBody
}

func newDockerClient() cli {
	nclient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli{c: nclient}
}

func Exec(req jkojsTypes.StartExecRequest, res *jkojsTypes.StartExecResponse) {
	c := newDockerClient()

	decodeSourceCode(&req)

	c.containerCreate(req)

	cfg := preparePacking(req)
	tarFile, err := packSourceAndCases(cfg)
	if err != nil {
		// ToDo: エラーログ
		lib.Logger.Sugar().Errorf("tarファイルに纏められませんでした: %s", err.Error())
		return
	}

	err = c.containerStart(res, tarFile)
	if err != nil {
		lib.Logger.Sugar().Errorf("コンテナ起動に失敗: %v", err.Error())
		return
	}
}

/*
	方針:
	1. コードをデコード ok
	2. コンテナ作成 ok
	3. コンテナに送るためのファイル類の準備 ok
	4. tarにまとめる ok
	5. コンテナに送る ok
	6. コンテナを起動 ok
	7. コンテナのログ取る -> ログからは実行結果取らない ok
	8. コンテナからワーカーが吐いたファイルを引っ張ってくる ok
	9. 終わったらコンテナを削除 ok
*/

func (dCli *cli) containerCreate(arg jkojsTypes.StartExecRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	f := false
	swappness := int64(0)
	PidsLimit := int64(512)

	res, err := dCli.c.ContainerCreate(ctx, &container.Config{
		Image:           "24abd",
		NetworkDisabled: true, // ネットワークを切る
		// Cmd:             []string{"tail", "-f", "/dev/null"},
		Cmd: []string{"/jkworker", "-lang", arg.Lang, "-id", arg.ProblemID}, // 実行する時のコマンド
		Tty: false,                                                          // Falseにしておく
	}, &container.HostConfig{
		AutoRemove:  false,  // これをtrueにすると実行結果が取れなくなる
		NetworkMode: "none", // ネットワークにつながらないようにする
		Resources: container.Resources{
			Memory: func() int64 {
				MaxMemorySize := "1024M" // メモリ制限 コンテナ1つ1024メガバイトまで
				mem, _ := strconv.ParseInt(MaxMemorySize, 10, 64)
				return mem
			}(),
			MemorySwap: func() int64 {
				MaxMemorySize := "0M" // スワップは0にする
				mem, _ := strconv.ParseInt(MaxMemorySize, 10, 64)
				return mem
			}(),
			OomKillDisable:   &f,         // メモリが溢れたときに自動ストップをかけておく
			MemorySwappiness: &swappness, // スワップを切る
			PidsLimit:        &PidsLimit, // フォーク爆弾を防ぐために低く設定しておく
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}
	dCli.container = res

	return nil
}

func (dCli cli) containerStart(arg *jkojsTypes.StartExecResponse, codes bytes.Buffer) error {
	// コンテナにファイルを送る
	err := dCli.c.CopyToContainer(context.Background(), dCli.container.ID, "/", &codes, types.CopyToContainerOptions{})
	if err != nil {
		lib.Logger.Sugar().Errorf("コンテナにファイルを送れませんでした: %v", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// ToDo: ここのオプションをちゃんと指定する
	if err := dCli.c.ContainerStart(ctx, dCli.container.ID, types.ContainerStartOptions{}); err != nil {
		lib.Logger.Sugar().Errorf("コンテナの起動に失敗しました: %v", err.Error())
		panic(err)
	}
	defer func() {
		err = dCli.c.ContainerRemove(ctx, dCli.container.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		})
		if err != nil {
			return
			// ToDo: ERR-LOG吐く
		}
	}()

	statusCh, errCh := dCli.c.ContainerWait(ctx, dCli.container.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			lib.Logger.Sugar().Errorf("コンテナ実行中にエラーが発生しました: %v", err.Error())
			return err
		}
	case <-statusCh:
	}

	// workerの実行結果を取ってくる
	f, _, err := dCli.c.CopyFromContainer(ctx, dCli.container.ID, "/out.json")
	if err != nil {
		panic(err)
	}
	b, _ := io.ReadAll(f)
	// 要らないデータがあるので取り除く
	to := trimer(b)
	os.WriteFile("test.json", to, 0660)

	err = json.Unmarshal(to, arg)
	if err != nil {
		lib.Logger.Sugar().Errorf("JSONのパースに失敗: %v", err.Error())
		return err
	}

	return nil
}
