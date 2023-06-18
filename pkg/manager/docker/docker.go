package docker

import (
	"archive/tar"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mct-joken/jkojs-agent/pkg/lib"
	"github.com/mct-joken/jkojs-agent/pkg/manager"
	"io"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type WorkerManager struct {
	cli *client.Client
}

func NewWorkerManager(c *client.Client) *WorkerManager {
	return &WorkerManager{cli: c}
}

func (m WorkerManager) Start(ctx context.Context, args manager.StartWorkerArgs) (manager.WorkerResponse, error) {
	tarFile, err := m.prepareFiles(args)
	if err != nil {
	}

	containerID, err := m.createContainer(ctx, args)
	if err != nil {
		fmt.Println(err)
		return manager.WorkerResponse{}, err
	}
	res, err := m.startContainer(tarFile, containerID)
	if err != nil {
		lib.Logger.Sugar().Errorf("コンテナ起動に失敗: %v", err)
		return manager.WorkerResponse{}, err
	}
	return res, nil
}

func (m WorkerManager) prepareFiles(req manager.StartWorkerArgs) (io.Reader, error) {
	err := decodeSourceCode(&req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cfg := preparePacking(req)
	tarFile, err := packSourceAndCases(cfg)
	if err != nil {
		lib.Logger.Sugar().Errorf("tarファイルに纏められませんでした: %s", err.Error())
		return nil, err
	}
	return &tarFile, nil
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

// createContainer Dockerコンテナを作成
func (m WorkerManager) createContainer(ctx context.Context, arg manager.StartWorkerArgs) (string, error) {
	f := false
	swappness := int64(0)
	PidsLimit := int64(512)
	//fmt.Println([]string{"/jkworker", "-lang", arg.Lang, "-id", arg.ProblemID})
	if lib.Config.ID == "" {
		fmt.Println("no image id")
		return "", errors.New("image id is not found")
	}

	command := []string{"/home/worker/ojs-worker", "-lang", arg.Lang, "-id", arg.ProblemID, "-p"}
	res, err := m.cli.ContainerCreate(ctx, &container.Config{
		Image:           lib.Config.ID,
		NetworkDisabled: true,    // ネットワークを切る
		Cmd:             command, // 実行する時のコマンド
		Tty:             false,   // Falseにしておく
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
		fmt.Println(err)
		return "", err
	}
	return res.ID, nil
}

// startContainer コンテナを起動
func (m WorkerManager) startContainer(codes io.Reader, containerID string) (manager.WorkerResponse, error) {
	// コンテナにファイルを送る
	err := m.cli.CopyToContainer(context.Background(), containerID, "/home/worker", codes, types.CopyToContainerOptions{})
	if err != nil {
		lib.Logger.Sugar().Errorf("コンテナにファイルを送れませんでした: %v", err.Error())
		return manager.WorkerResponse{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// ToDo: ここのオプションをちゃんと指定する
	err = m.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		lib.Logger.Sugar().Errorf("コンテナの起動に失敗しました: %v", err)
		return manager.WorkerResponse{}, err
	}
	//defer func() {
	//	err = dCli.c.ContainerRemove(ctx, dCli.container.ID, types.ContainerRemoveOptions{
	//		RemoveVolumes: true,
	//		RemoveLinks:   false,
	//		Force:         true,
	//	})
	//	if err != nil {
	//		return
	//		// ToDo: ERR-LOG吐く
	//	}
	//}()

	statusCh, errCh := m.cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			lib.Logger.Sugar().Errorf("コンテナ実行中にエラーが発生しました: %v", err.Error())
			return manager.WorkerResponse{}, err
		}
	case <-statusCh:
	}

	// workerの実行結果を取ってくる
	f, _, err := m.cli.CopyFromContainer(ctx, containerID, "/home/worker/out.json")
	if err != nil {
		panic(err)
	}

	var out []byte
	reader := tar.NewReader(f)
	for {
		h, err := reader.Next()
		if err == io.EOF {
			break
		}

		if h.Name == "out.json" {
			out, _ = io.ReadAll(reader)
			break
		}
	}

	res := manager.WorkerResponse{}
	err = json.Unmarshal(out, &res)
	if err != nil {
		lib.Logger.Sugar().Errorf("JSONのパースに失敗: %v", err.Error())
		return manager.WorkerResponse{}, err
	}

	return res, nil
}
