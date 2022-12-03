package docker

import (
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func Exec(arg jkojsTypes.StartExecRequest) {
	c := newDockerClient()

	decodeSourceCode(&arg)
	c.containerCreate()
	c.containerStart(arg)
}

/*
	方針:
	1. コードをデコード ok
	2. コンテナ作成 ok
	3. コンテナに送るためのファイル類の準備 todo
	4. tarにまとめる todo
	5. コンテナに送る todo
	6. コンテナを起動 ok
	7. コンテナのログ取る -> ログからは実行結果取らない ok
	8. コンテナからワーカーが吐いたファイルを引っ張ってくる ok
	9. 終わったらコンテナを削除 ok
*/

func (dCli *cli) containerCreate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	f := false
	swappness := int64(0)
	PidsLimit := int64(512)

	res, err := dCli.c.ContainerCreate(ctx, &container.Config{
		Image:           "1e",
		NetworkDisabled: true,                                                       // ネットワークを切る
		Cmd:             []string{"/jkworker", "-lang", "Clang++", "-id", "112233"}, // 実行する時のコマンド
		Tty:             false,                                                      // Falseにしておく
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

func (dCli cli) containerStart(arg jkojsTypes.StartExecRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// ToDo: ここのオプションをちゃんと指定する
	if err := dCli.c.ContainerStart(ctx, dCli.container.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := dCli.c.ContainerWait(ctx, dCli.container.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
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

	_ = dCli.c.ContainerRemove(ctx, dCli.container.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	})

	return nil
}
