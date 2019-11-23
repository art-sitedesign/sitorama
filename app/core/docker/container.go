package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

// CreateContainer создаст контейнер
func (d *Docker) CreateContainer(ctx context.Context, name string, config *container.Config, hostConfig *container.HostConfig) (string, error) {
	networkingConfig := &network.NetworkingConfig{}

	res, err := d.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, containerName(name))
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
	args.Add("name", containerName(name))

	opts := types.ContainerListOptions{
		Quiet:   false,
		Size:    false,
		All:     false,
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

func containerName(s string) string {
	return fmt.Sprintf("%s-%s", prefix, s)
}
