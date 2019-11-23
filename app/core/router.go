package core

import (
	"context"
	"pet-projects/sitorama/app/core/docker"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

const RouterName = "router"

func (c *Core) findRouter(ctx context.Context) (*types.Container, error) {
	containers, err := c.docker.FindContainers(ctx, RouterName)
	if err != nil {
		return nil, errors.Wrap(err, "failed get router")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func (c *Core) createRouter(ctx context.Context) (string, error) {
	portSet, portMap := docker.BindPorts(map[string]string{"80": "80"})

	config := docker.DefaultContainerConfig()
	config.ExposedPorts = portSet
	config.Image = "nginx:latest"

	hostConfig := docker.DefaultContainerHostConfig()
	hostConfig.PortBindings = portMap

	return c.docker.CreateContainer(ctx, RouterName, config, hostConfig)
}
