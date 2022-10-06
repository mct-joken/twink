package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strconv"
	"time"
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

func (wo *WorkSpace) Create(containerName string, imageName string, port string) (string, error) {
	// ポート番号を10000~65535に限定する
	err := func(p string) error {
		i, e := strconv.Atoi(p)
		if e != nil {
			return e
		}

		if i < 65535 || i >= 10000 {
			return nil
		}
		return errors.New("invalid port number")
	}(port)
	if err != nil {
		return "", err
	}

	// containerConfig コンテナの設定
	containerConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{nat.Port("22"): struct{}{}},
		Tty:          true, // ここがfalseだと、プロセスが終了するとコンテナが終了する
	}

	// containerHostConfig コンテナを実行するサーバーのリソースなど
	containerHostConfig := &container.HostConfig{
		AutoRemove: false, // これをtrueにするとコンテナが停止するとコンテナが削除されてしまう
		Resources: container.Resources{
			OomKillDisable: func() *bool { b := false; return &b }(), //これをtrueにすると、ホスト側のメモリを食いつぶすことがあるのでfalseにする
		},
		PortBindings: map[nat.Port][]nat.PortBinding{
			nat.Port("22"): {{HostPort: port}},
		},
	}

	// コンテナを作る
	co, err := cli.ContainerCreate(context.Background(), containerConfig, containerHostConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}

	// コンテナIDを返す
	return co.ID, err
}

func (wo *WorkSpace) Start(containerID string) error {
	// ToDo: Docker コンテナをスタート
	err := cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (wo *WorkSpace) Stop(containerID string) error {
	// ToDo: Docker コンテナを停止
	// コンテナ終了のシグナルを送っても停止しない場合は、10秒後に強制停止する
	err := cli.ContainerStop(context.Background(), containerID, func() *time.Duration { t := 1 * time.Second; return &t }())
	if err != nil {
		return err
	}
	return nil
}

func (wo *WorkSpace) Destroy() {
	// ToDo: Docker コンテナを削除
}
