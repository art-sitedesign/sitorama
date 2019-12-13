package builder

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

const (
	pgDefaultUser = "postgres"
	pgDataDir     = "postgres"

	postgresDefaultPort = 5432
	postgresForwardPort = "forwardPort"
)

type Postgres struct {
	docker *docker.Docker
	name   string
	config Config
}

func NewPostgres(docker *docker.Docker, name string) Builder {
	return &Postgres{
		docker: docker,
		name:   name,
	}
}

func (p *Postgres) Name() string {
	return "Postgres"
}

func (p *Postgres) ConfigNames() []string {
	return []string{
		postgresForwardPort,
	}
}

func (p *Postgres) ConfigByName(name string) (string, error) {
	switch name {
	case postgresForwardPort:
		freePort := utils.FindNearPort(postgresDefaultPort)
		return strconv.Itoa(freePort), nil
	default:
		return "", errors.New("unknown config " + name)
	}
}

func (p *Postgres) SetConfig(config Config) {
	p.config = config
}

func (p *Postgres) Checker() (string, error) {
	return renderChecker("postgres.php", map[string]string{
		"DBHost": p.alias(),
		"DBName": p.dbName(),
		"DBUser": p.user(),
		"DBPass": p.password(),
	})
}

func (p *Postgres) Info() map[string]string {
	return map[string]string{
		"Public host":  "127.0.0.1",
		"Public port":  p.portToForward(),
		"Private host": p.alias(),
		"Private port": services.PostgresPort,
		"User":         p.user(),
		"Password":     p.password(),
		"DB name":      p.dbName(),
	}
}

func (p *Postgres) Build(ctx context.Context) error {
	// поиск сети приложения
	network, err := p.docker.FindNetwork(ctx)
	if err != nil {
		return err
	}

	// сборка контейнера apache
	err = p.buildContainer(ctx, network.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) buildContainer(ctx context.Context, networkID string) error {
	pgDataPath, err := utils.ProjectVolumeFullPath(p.name, pgDataDir)
	if err != nil {
		return err
	}

	servicePostgres := services.NewPostgres(
		p.docker,
		p.name,
		p.portToForward(),
		p.user(),
		p.password(),
		p.dbName(),
		pgDataPath,
	)

	container, err := servicePostgres.Find(ctx)
	if err != nil {
		return err
	}

	if container == nil {
		// если контейнер не найден - создаем
		cID, err := servicePostgres.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = p.docker.ConnectNetwork(ctx, networkID, cID, []string{p.alias()})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = p.docker.StartContainer(ctx, cID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) alias() string {
	return fmt.Sprintf("%s.postgres", p.name)
}

func (p *Postgres) user() string {
	return pgDefaultUser
}

func (p *Postgres) password() string {
	return defaultPassword
}

func (p Postgres) dbName() string {
	return strings.Replace(p.name, ".", "_", -1)
}

func (p *Postgres) portToForward() string {
	confPort := p.config.String(postgresForwardPort)
	if confPort != nil {
		return *confPort
	}

	return strconv.Itoa(postgresDefaultPort)
}
