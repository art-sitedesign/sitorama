package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

const (
	defaultUser     = "postgres"
	defaultPassword = "sitorama"
	pgDataDir       = "postgres"
)

type Postgres struct {
	docker *docker.Docker
	name   string
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
	return []string{}
}

func (p *Postgres) ConfigByName(name string) (string, error) {
	return "", nil
}

func (p *Postgres) SetConfig(config Config) {
}

func (p *Postgres) Checker() (string, error) {
	return renderChecker("postgres.php", map[string]string{
		"DBHost": p.alias(),
		"DBName": p.dbName(),
		"DBUser": p.user(),
		"DBPass": p.password(),
	})
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

	servicePostgres := services.NewPostgres(p.docker, p.name, p.user(), p.password(), p.dbName(), pgDataPath)

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
	return defaultUser
}

func (p *Postgres) password() string {
	return defaultPassword
}

func (p Postgres) dbName() string {
	return strings.Replace(p.name, ".", "_", -1)
}
