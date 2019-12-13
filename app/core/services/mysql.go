package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
)

const (
	MySQLPort     = "3306"
	mysqlDataPath = "/var/lib/mysql"
)

type MySQL struct {
	docker      *docker.Docker
	name        string
	forwardPort string
	dbPass      string
	dbName      string
	dataPath    string
}

func NewMySQL(
	docker *docker.Docker,
	name string,
	forwardPort string,
	dbPass string,
	dbName string,
	dataPath string,
) *MySQL {
	return &MySQL{
		docker:      docker,
		name:        name,
		forwardPort: forwardPort,
		dbPass:      dbPass,
		dbName:      dbName,
		dataPath:    dataPath,
	}
}

func (m *MySQL) Find(ctx context.Context) (*types.Container, error) {
	containers, err := m.docker.FindContainers(ctx, m.ContainerName())
	if err != nil {
		return nil, errors.Wrap(err, "failed get container mysql")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func (m *MySQL) Create(ctx context.Context) (string, error) {
	portSet, portMap := docker.BindPorts(map[string]string{m.forwardPort: MySQLPort})

	config := docker.DefaultContainerConfig()
	config.ExposedPorts = portSet
	config.Image = "library/mysql:8"
	config.Env = []string{
		fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", m.dbPass),
		fmt.Sprintf("MYSQL_DATABASE=%s", m.dbName),
	}

	hostConfig := docker.DefaultContainerHostConfig()
	hostConfig.PortBindings = portMap

	volumes := docker.MakeVolumes(map[string]string{m.dataPath: mysqlDataPath})
	hostConfig.Mounts = volumes

	cID, err := m.docker.CreateContainer(ctx, m.ContainerName(), config, hostConfig)
	if err != nil {
		return "", err
	}

	return cID, nil
}

func (m *MySQL) ContainerName() string {
	return fmt.Sprintf("%s_mysql", m.name)
}
