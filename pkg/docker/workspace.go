package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type WorkSpace struct{}

var cli *client.Client

func NewConnection() {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli = c
	if err != nil {
		return
	}
}

func (wo *WorkSpace) Create(imageName string) (string, error) {
	// FIXME: 外部からイメージの指定できるようにする
	// containerConfig コンテナの設定
	containerConfig := &container.Config{
		Image: imageName,
		Tty:   true, // ここがfalseだと、プロセスが終了するとコンテナが終了する
	}

	// containerHostConfig コンテナを実行するサーバーのリソースなど
	containerHostConfig := &container.HostConfig{
		AutoRemove: false, // これをtrueにするとコンテナが停止するとコンテナが削除されてしまう
		Resources: container.Resources{
			OomKillDisable: func() *bool { b := false; return &b }(),
		},
	}

	// コンテナを作る
	co, err := cli.ContainerCreate(context.Background(), containerConfig, containerHostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	// コンテナIDを返す
	return co.ID, err
}

func (wo *WorkSpace) Start() {
	// ToDo: Docker コンテナをスタート
}

func (wo *WorkSpace) Stop() {
	// ToDo: Docker コンテナを停止
}

func (wo *WorkSpace) Destroy() {
	// ToDo: Docker コンテナを削除
}
