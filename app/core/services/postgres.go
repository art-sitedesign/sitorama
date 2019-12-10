package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
)

const (
	defaultPGData = "/var/lib/postgresql/data/pgdata"
)

type Postgres struct {
	docker      *docker.Docker
	name        string
	forwardPort string
	dbUser      string
	dbPass      string
	dbName      string
	pgDataPath  string
}

func NewPostgres(
	docker *docker.Docker,
	name string,
	forwardPort string,
	dbUser string,
	dbPass string,
	dbName string,
	pgDataPath string,
) *Postgres {
	return &Postgres{
		docker:      docker,
		name:        name,
		forwardPort: forwardPort,
		dbUser:      dbUser,
		dbPass:      dbPass,
		dbName:      dbName,
		pgDataPath:  pgDataPath,
	}
}

func (p *Postgres) Find(ctx context.Context) (*types.Container, error) {
	containers, err := p.docker.FindContainers(ctx, p.ContainerName())
	if err != nil {
		return nil, errors.Wrap(err, "failed get container postgres")
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &containers[0], nil
}

func (p *Postgres) Create(ctx context.Context) (string, error) {
	portSet, portMap := docker.BindPorts(map[string]string{p.forwardPort: "5432"})

	config := docker.DefaultContainerConfig()
	config.ExposedPorts = portSet
	config.Image = "library/postgres:12"
	config.Env = []string{
		fmt.Sprintf("POSTGRES_USER=%s", p.dbUser),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", p.dbPass),
		fmt.Sprintf("POSTGRES_DB=%s", p.dbName),
		fmt.Sprintf("PGDATA=%s", defaultPGData),
	}

	hostConfig := docker.DefaultContainerHostConfig()
	hostConfig.PortBindings = portMap

	volumes := docker.MakeVolumes(map[string]string{p.pgDataPath: defaultPGData})
	hostConfig.Mounts = volumes

	cID, err := p.docker.CreateContainer(ctx, p.ContainerName(), config, hostConfig)
	if err != nil {
		return "", err
	}

	return cID, nil
}

func (p *Postgres) ContainerName() string {
	return fmt.Sprintf("%s_postgres", p.name)
}
