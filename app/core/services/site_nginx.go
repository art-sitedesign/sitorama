package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type SiteNginx struct {
	docker  *docker.Docker
	name    string
	pfAlias string
}

func NewSiteNginx(d *docker.Docker, n string, pfAlias string) *SiteNginx {
	return &SiteNginx{
		docker:  d,
		name:    n,
		pfAlias: pfAlias,
	}
}

func (sn *SiteNginx) Find(ctx context.Context) (*types.Container, error) {
	containers, err := sn.docker.FindContainers(ctx, sn.ContainerName())
	if err != nil {
		return nil, errors.Wrap(err, "failed get container site nginx")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func (sn *SiteNginx) Create(ctx context.Context) (string, error) {
	portSet, _ := docker.BindPorts(map[string]string{"80": "80"})

	config := docker.DefaultContainerConfig()
	config.ExposedPorts = portSet
	config.Image = "nginx:latest"

	hostConfig := docker.DefaultContainerHostConfig()

	cID, err := sn.docker.CreateContainer(ctx, sn.ContainerName(), config, hostConfig)
	if err != nil {
		return "", err
	}

	nginxConf, err := utils.RenderTemplateInBuffer(utils.SiteNginxBaseTemplate, nil)
	if err != nil {
		return "", err
	}

	params := map[string]string{"Domain": sn.name, "PFAlias": sn.pfAlias}
	serverConf, err := utils.RenderTemplateInBuffer(utils.SiteNginxServerTemplate, params)

	err = sn.docker.CopyToContainer(ctx, cID, "/etc/nginx/", "nginx.conf", nginxConf)
	if err != nil {
		return "", err
	}

	err = sn.docker.CopyToContainer(ctx, cID, "/etc/nginx/conf.d/", "default.conf", serverConf)
	if err != nil {
		return "", err
	}

	return cID, nil
}

func (sn *SiteNginx) ContainerName() string {
	return utils.ContainerName(fmt.Sprintf("%s_nginx", sn.name))
}
