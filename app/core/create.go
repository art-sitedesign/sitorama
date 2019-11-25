package core

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

func (c *Core) CreateSite(ctx context.Context, name string) error {
	pfAlias := fmt.Sprintf("%s.phpfpm", name)
	ngAlias := fmt.Sprintf("%s.nginx", name)

	// создание конфиг-фалйа для роутера
	err := createRouterConfig(name, ngAlias)
	if err != nil {
		return err
	}

	// поиск сети приложения
	network, err := c.docker.FindNetwork(ctx)
	if err != nil {
		return err
	}

	// создание PHP-FPM контейнера
	sitePHPFPM := services.NewSitePHPFPM(c.docker, name)
	pfContainer, err := sitePHPFPM.Find(ctx)
	if err != nil {
		return err
	}

	if pfContainer == nil {
		pfCID, err := sitePHPFPM.Create(ctx)
		if err != nil {
			return err
		}

		err = c.docker.ConnectNetwork(ctx, network.ID, pfCID, []string{pfAlias})
		if err != nil {
			return err
		}

		err = c.docker.StartContainer(ctx, pfCID)
	}

	// создание контейнера
	siteNginx := services.NewSiteNginx(c.docker, name, pfAlias)
	container, err := siteNginx.Find(ctx)
	if err != nil {
		return err
	}

	if container == nil {
		cID, err := siteNginx.Create(ctx)
		if err != nil {
			return err
		}

		err = c.docker.ConnectNetwork(ctx, network.ID, cID, []string{ngAlias})
		if err != nil {
			return err
		}

		err = c.docker.StartContainer(ctx, cID)
		if err != nil {
			return err
		}
	}

	// перезапуск роутера для применения конфига
	r := services.NewRouter(c.docker)
	router, err := r.Find(ctx)
	if err != nil {
		return err
	}

	err = c.docker.RestartContainer(ctx, router.ID)
	if err != nil {
		return err
	}

	return nil
}

func createRouterConfig(name string, containerAlias string) error {
	tmpl := template.Must(template.ParseFiles(utils.RouterConfTemplate))

	err := os.MkdirAll(utils.RouterConfDir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(utils.RouterConfPath(name))
	if err != nil {
		return err
	}
	defer f.Close()

	data := map[string]string{
		"Domain":         name,
		"ContainerAlias": containerAlias,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
