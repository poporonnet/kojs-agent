package docker

import (
	"context"
	"fmt"
	"io"
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
	1. コードをデコード
	2. コンテナ作成
	3. コンテナに送るためのファイル類の準備
	4. tarにまとめる
	5. コンテナに送る
	6. コンテナを起動
	7. コンテナのログ取る -> ログからは実行結果取らない
	8. コンテナからワーカーが吐いたファイルを引っ張ってくる
	9. 終わったらコンテナを削除
*/

func (dCli *cli) containerCreate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	f := false
	swappness := int64(0)   // スワップを封印
	PidsLimit := int64(512) // 対フォーク爆弾の設定
	res, err := dCli.c.ContainerCreate(ctx, &container.Config{
		Image:           "6f579486009b7e599f09285d40f427b6a7cc235cbd52037a548d91ba4b3c917c",
		NetworkDisabled: true,
		Cmd:             []string{"cat", "/etc/os-release"}, // 実行する時のコマンド
		Tty:             false,                              // Falseにしておく
	}, &container.HostConfig{
		AutoRemove:  false, // これをOnにするとLogが取れなくなって死ぬ
		NetworkMode: "none",
		Resources: container.Resources{
			Memory: func() int64 {
				MaxMemorySize := "512M" // メモリ制限 コンテナ1つ512メガバイトまで
				mem, _ := strconv.ParseInt(MaxMemorySize, 10, 64)
				return mem
			}(),
			MemorySwap: func() int64 {
				MaxMemorySize := "512M" // メモリ制限 コンテナ1つ512メガバイトまで
				mem, _ := strconv.ParseInt(MaxMemorySize, 10, 64)
				return mem
			}(),
			OomKillDisable:   &f,
			MemorySwappiness: &swappness,
			PidsLimit:        &PidsLimit,
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

	out, err := dCli.c.ContainerLogs(ctx, dCli.container.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	f, _, err := dCli.c.CopyFromContainer(ctx, dCli.container.ID, "/out.json")
	if err != nil {
		panic(err)
	}
	b, _ := io.ReadAll(f)
	fmt.Println(string(b))

	_ = dCli.c.ContainerRemove(ctx, dCli.container.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	})

	c, e := io.ReadAll(out)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(c))
	return nil
}
