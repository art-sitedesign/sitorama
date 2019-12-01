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

type SiteNginx struct {
	docker     *docker.Docker
	name       string
	entryPoint string
	pfAlias    string
	config     *string
}

func NewSiteNginx(d *docker.Docker, n string, ep string, pfAlias string, cf *string) *SiteNginx {
	return &SiteNginx{
		docker:     d,
		name:       n,
		entryPoint: ep,
		pfAlias:    pfAlias,
		config:     cf,
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

	var serverConf *bytes.Buffer

	if sn.config == nil {
		// если конфиг в конструктор не был передан рендерим дефолтный
		serverConf, err = sn.RenderConfig()
		if err != nil {
			return "", err
		}
	} else {
		// если был передан, используем его
		b := bytes.Buffer{}
		b.WriteString(*sn.config)
		serverConf = &b
	}

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

func (sn *SiteNginx) RenderConfig() (*bytes.Buffer, error) {
	params := map[string]string{"Domain": sn.name, "EntryPoint": sn.entryPoint, "PFAlias": sn.pfAlias}
	return utils.RenderTemplateInBuffer(utils.SiteNginxServerTemplate, params)
}

func (sn *SiteNginx) ContainerName() string {
	return fmt.Sprintf("%s_nginx", sn.name)
}
