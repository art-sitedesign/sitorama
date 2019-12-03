package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type Apache struct {
	docker      *docker.Docker
	name        string
	projectPath string
	entryPoint  string
	config      *string
}

func NewApache(docker *docker.Docker, name string, projectPath string, entryPoint string, config *string) *Apache {
	return &Apache{
		docker:      docker,
		name:        name,
		projectPath: projectPath,
		entryPoint:  entryPoint,
		config:      config,
	}
}

func (a *Apache) Find(ctx context.Context) (*types.Container, error) {
	containers, err := a.docker.FindContainers(ctx, a.ContainerName())
	if err != nil {
		return nil, errors.Wrap(err, "failed get container apache")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func (a *Apache) Create(ctx context.Context) (string, error) {
	portSet, _ := docker.BindPorts(map[string]string{"80": "80"})

	config := docker.DefaultContainerConfig()
	//config.User = "www-data:www-data"
	//config.WorkingDir = "/app"
	config.ExposedPorts = portSet
	config.Image = "library/php:7.2-apache"

	hostConfig := docker.DefaultContainerHostConfig()

	volumes := docker.MakeVolumes(map[string]string{a.projectPath: "/var/www/app"})
	hostConfig.Mounts = volumes

	cID, err := a.docker.CreateContainer(ctx, a.ContainerName(), config, hostConfig)
	if err != nil {
		return "", err
	}

	var serverConf *bytes.Buffer

	if a.config == nil {
		// если конфиг в конструктор не был передан рендерим дефолтный
		serverConf, err = a.RenderConfig()
		if err != nil {
			return "", err
		}
	} else {
		// если был передан, используем его
		b := bytes.Buffer{}
		b.WriteString(*a.config)
		serverConf = &b
	}

	err = a.docker.CopyToContainer(ctx, cID, "/etc/apache2/sites-available/", "000-default.conf", serverConf)
	if err != nil {
		return "", err
	}

	return cID, nil
}

func (a *Apache) ContainerName() string {
	return fmt.Sprintf("%s_apache", a.name)
}

func (a *Apache) RenderConfig() (*bytes.Buffer, error) {
	params := map[string]string{"Domain": a.name, "EntryPoint": a.entryPoint}
	return utils.RenderTemplateInBuffer(utils.ApacheServerTemplate, params)
}
