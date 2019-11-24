package services

import (
	"context"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type Router struct {
	docker *docker.Docker
}

func NewRouter(d *docker.Docker) Service {
	return &Router{docker: d}
}

// Find найдет контейнер роутера
func (r *Router) Find(ctx context.Context) (*types.Container, error) {
	containers, err := r.docker.FindContainers(ctx, utils.RouterName)
	if err != nil {
		return nil, errors.Wrap(err, "failed get router")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

// Create создаст контейнер роутера
func (r *Router) Create(ctx context.Context) (string, error) {
	portSet, portMap := docker.BindPorts(map[string]string{"80": "80"})

	config := docker.DefaultContainerConfig()
	config.ExposedPorts = portSet
	config.Image = "nginx:latest"

	hostConfig := docker.DefaultContainerHostConfig()
	hostConfig.PortBindings = portMap

	err := os.MkdirAll(utils.RouterConfDir, 0755)
	if err != nil {
		return "", err
	}

	path, err := filepath.Abs(utils.RouterConfDir)
	if err != nil {
		return "", err
	}

	volumes := docker.MakeVolumes(map[string]string{path: "/etc/nginx/conf.d/"})
	hostConfig.Mounts = volumes

	return r.docker.CreateContainer(ctx, utils.RouterName, config, hostConfig)
}
