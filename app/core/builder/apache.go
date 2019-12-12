package builder

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

const (
	apacheServerConfig = "apache-server-conf"
)

type Apache struct {
	docker     *docker.Docker
	name       string
	entryPoint string
	config     Config
}

func NewApache(docker *docker.Docker, name string, entryPoint string) Builder {
	return &Apache{
		docker:     docker,
		name:       name,
		entryPoint: entryPoint,
	}
}

func (a *Apache) Name() string {
	return "Apache"
}

func (a *Apache) ConfigNames() []string {
	return []string{
		apacheServerConfig,
	}
}

func (a *Apache) ConfigByName(name string) (string, error) {
	switch name {
	case apacheServerConfig:
		conf, err := a.serverConfig()
		if err != nil {
			return "", err
		}
		return conf, nil
	default:
		return "", errors.New("unknown config " + name)
	}
}

func (a *Apache) SetConfig(config Config) {
	a.config = config
}

func (a *Apache) Checker() (string, error) {
	return renderChecker("base.php", map[string]string{"Name": a.name})
}

func (a *Apache) Info() map[string]string {
	return map[string]string{
		"Private host": a.alias(),
	}
}

func (a *Apache) Build(ctx context.Context) error {
	// создание конфиг-фалйа для роутера
	err := utils.CreateRouterConfig(a.name, a.alias())
	if err != nil {
		return err
	}

	// поиск сети приложения
	network, err := a.docker.FindNetwork(ctx)
	if err != nil {
		return err
	}

	// сборка контейнера apache
	err = a.buildContainer(ctx, network.ID)
	if err != nil {
		return err
	}

	// перезапуск роутера для применения конфига
	r := services.NewRouter(a.docker)
	router, err := r.Find(ctx)
	if err != nil {
		return err
	}

	err = a.docker.RestartContainer(ctx, router.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *Apache) buildContainer(ctx context.Context, networkID string) error {
	projectPath, err := utils.ProjectFullPath(a.name)
	conf := a.config.String(apacheServerConfig)

	serviceApache := services.NewApache(a.docker, a.name, projectPath, a.entryPoint, conf)

	container, err := serviceApache.Find(ctx)
	if err != nil {
		return err
	}

	if container == nil {
		// если контейнер не найден - создаем
		cID, err := serviceApache.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = a.docker.ConnectNetwork(ctx, networkID, cID, []string{a.alias()})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = a.docker.StartContainer(ctx, cID)
		if err != nil {
			return err
		}

		err = a.docker.ExecInContainer(ctx, cID, []string{
			"apt update",
			"apt -y install libpq-dev",
			"docker-php-ext-install pdo",
			"docker-php-ext-install pdo_pgsql",
			"docker-php-ext-install pgsql",
		})
		if err != nil {
			return err
		}

		err = a.docker.RestartContainer(ctx, cID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Apache) alias() string {
	return fmt.Sprintf("%s.apache", a.name)
}

func (a *Apache) serverConfig() (string, error) {
	projectPath, err := utils.ProjectFullPath(a.name)
	if err != nil {
		return "", err
	}

	serviceApache := services.NewApache(a.docker, a.name, projectPath, a.entryPoint, nil)
	conf, err := serviceApache.RenderConfig()
	if err != nil {
		return "", err
	}

	return conf.String(), nil
}
