package builder

import (
	"context"
	"fmt"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type NginxPHPFPM struct {
	docker *docker.Docker
	name   string
}

func NewNginxPHPFPM(d *docker.Docker, n string) Builder {
	return &NginxPHPFPM{
		docker: d,
		name:   n,
	}
}

func (np *NginxPHPFPM) Build(ctx context.Context) error {
	pfAlias := fmt.Sprintf("%s.phpfpm", np.name)
	ngAlias := fmt.Sprintf("%s.nginx", np.name)

	// создание конфиг-фалйа для роутера
	err := utils.CreateRouterConfig(np.name, ngAlias)
	if err != nil {
		return err
	}

	// поиск сети приложения
	network, err := np.docker.FindNetwork(ctx)
	if err != nil {
		return err
	}

	// сборка контейнера PHP-FPM
	err = np.buildPHPFPM(ctx, network.ID, pfAlias)
	if err != nil {
		return err
	}

	err = np.buildNginx(ctx, network.ID, pfAlias, ngAlias)
	if err != nil {
		return err
	}

	// перезапуск роутера для применения конфига
	r := services.NewRouter(np.docker)
	router, err := r.Find(ctx)
	if err != nil {
		return err
	}

	err = np.docker.RestartContainer(ctx, router.ID)
	if err != nil {
		return err
	}

	return nil
}

func (np *NginxPHPFPM) buildPHPFPM(ctx context.Context, networkID string, pfAlias string) error {
	// поиск PHP-FPM контейнера
	sitePHPFPM := services.NewSitePHPFPM(np.docker, np.name)
	pfContainer, err := sitePHPFPM.Find(ctx)
	if err != nil {
		return err
	}

	if pfContainer == nil {
		// если контейнер не найден - создаем
		pfCID, err := sitePHPFPM.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = np.docker.ConnectNetwork(ctx, networkID, pfCID, []string{pfAlias})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = np.docker.StartContainer(ctx, pfCID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (np *NginxPHPFPM) buildNginx(ctx context.Context, networkID string, pfAlias string, ngAlias string) error {
	siteNginx := services.NewSiteNginx(np.docker, np.name, pfAlias)
	container, err := siteNginx.Find(ctx)
	if err != nil {
		return err
	}

	if container == nil {
		// если контейнер не найден - создаем
		cID, err := siteNginx.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = np.docker.ConnectNetwork(ctx, networkID, cID, []string{ngAlias})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = np.docker.StartContainer(ctx, cID)
		if err != nil {
			return err
		}
	}

	return nil
}
