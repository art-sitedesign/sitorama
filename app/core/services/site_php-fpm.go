package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type SitePHPFPM struct {
	docker      *docker.Docker
	name        string
	projectPath string
}

func NewSitePHPFPM(docker *docker.Docker, name string, projectPath string) *SitePHPFPM {
	return &SitePHPFPM{
		docker:      docker,
		name:        name,
		projectPath: projectPath,
	}
}

func (sp *SitePHPFPM) Find(ctx context.Context) (*types.Container, error) {
	containers, err := sp.docker.FindContainers(ctx, sp.ContainerName())
	if err != nil {
		return nil, errors.Wrap(err, "failed get container site nginx")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func (sp *SitePHPFPM) Create(ctx context.Context) (string, error) {
	//portSet, _ := docker.BindPorts(map[string]string{"9000": "9000"})

	config := docker.DefaultContainerConfig()
	//config.ExposedPorts = portSet
	config.Image = "bitnami/php-fpm:latest"

	hostConfig := docker.DefaultContainerHostConfig()

	volumes := docker.MakeVolumes(map[string]string{sp.projectPath: "/app"})
	hostConfig.Mounts = volumes

	cID, err := sp.docker.CreateContainer(ctx, sp.ContainerName(), config, hostConfig)
	if err != nil {
		return "", err
	}

	phpIni, err := utils.RenderTemplateInBuffer(utils.PHPIniTemplate, nil)
	if err != nil {
		return "", err
	}

	phpLibs, err := utils.RenderTemplateInBuffer(utils.PHPLibsTemplate, nil)
	if err != nil {
		return "", err
	}

	err = sp.docker.CopyToContainer(ctx, cID, "/opt/bitnami/php/etc/conf.d/", "sitorama.ini", phpIni)
	if err != nil {
		return "", err
	}

	err = sp.docker.CopyToContainer(ctx, cID, "/opt/bitnami/php/etc/conf.d/", "libs.ini", phpLibs)
	if err != nil {
		return "", err
	}

	return cID, nil
}

func (sp *SitePHPFPM) ContainerName() string {
	return fmt.Sprintf("%s_php-fpm", sp.name)
}
