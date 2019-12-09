package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"github.com/art-sitedesign/sitorama/app/utils"
)

// CreateContainer создаст контейнер
func (d *Docker) CreateContainer(ctx context.Context, name string, config *container.Config, hostConfig *container.HostConfig) (string, error) {
	networkingConfig := &network.NetworkingConfig{}

	images, err := d.FindImage(ctx, config.Image)
	if err != nil {
		return "", nil
	}

	if len(images) == 0 {
		err = d.PullImage(ctx, config.Image)
		if err != nil {
			return "", err
		}
	}

	res, err := d.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, utils.ContainerName(name))
	if err != nil {
		return "", err
	}

	//todo: log warnings
	fmt.Println(res.Warnings)

	return res.ID, nil
}

// FindContainers найдет контейнеры по названию
func (d *Docker) FindContainers(ctx context.Context, name string) ([]types.Container, error) {
	args := filters.NewArgs()
	args.Add("name", utils.ContainerName(name))

	opts := types.ContainerListOptions{
		Quiet:   false,
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: args,
	}

	return d.client.ContainerList(ctx, opts)
}

// RunContainer запустит контейнер
func (d *Docker) StartContainer(ctx context.Context, containerID string) error {
	options := types.ContainerStartOptions{
		CheckpointID:  "",
		CheckpointDir: "",
	}

	return d.client.ContainerStart(ctx, containerID, options)
}

// RestartContainer перезапустит контейнер
func (d *Docker) RestartContainer(ctx context.Context, containerID string) error {
	return d.client.ContainerRestart(ctx, containerID, nil)
}

// StopContainer остановит контейнер
func (d *Docker) StopContainer(ctx context.Context, containerID string) error {
	return d.client.ContainerStop(ctx, containerID, nil)
}

// RemoveContainer удалит контейнер
func (d *Docker) RemoveContainer(ctx context.Context, containerID string) error {
	options := types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         false,
	}
	return d.client.ContainerRemove(ctx, containerID, options)
}

// CopyToContainer скопирует данные в файловую систему контейнера
func (d *Docker) CopyToContainer(ctx context.Context, containerID string, path string, fileName string, data *bytes.Buffer) error {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: fileName,
		Mode: 0755,
		Size: int64(data.Len()),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := tw.Write(data.Bytes()); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	options := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	}

	return d.client.CopyToContainer(ctx, containerID, path, &buf, options)
}

// ExecInContainer выполнит bash команду внутри контейнера
func (d *Docker) ExecInContainer(ctx context.Context, containerID string, commands []string) error {
	for _, command := range commands {
		config := types.ExecConfig{
			User:         "",
			Privileged:   false,
			Tty:          false,
			AttachStdin:  false,
			AttachStderr: false,
			AttachStdout: false,
			Detach:       false,
			DetachKeys:   "",
			Env:          nil,
			Cmd:          []string{"/bin/bash", "-c", command},
		}
		resp, err := d.client.ContainerExecCreate(ctx, containerID, config)
		if err != nil {
			return err
		}

		startConfig := types.ExecStartCheck{
			Detach: false,
			Tty:    false,
		}
		err = d.client.ContainerExecStart(ctx, resp.ID, startConfig)
		if err != nil {
			return err
		}

		for i := 0; i < 60; i++ {
			// делаем 60 попыток проверить статус команды и если нет - выходим (за минуту нужно успеть)
			inspectResp, err := d.client.ContainerExecInspect(ctx, resp.ID)
			if err != nil {
				return err
			}

			if !inspectResp.Running {
				break
			}

			time.Sleep(time.Millisecond * 1000)
		}
	}

	return nil
}
