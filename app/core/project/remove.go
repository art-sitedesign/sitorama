package project

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/filesystem"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/core/settings"
	"github.com/art-sitedesign/sitorama/app/utils"
)

func (p *Project) Remove(ctx context.Context, name string) error {
	containers, err := p.docker.FindContainers(ctx, name)
	if err != nil {
		return err
	}

	// удаление контейнеров проекта
	for _, container := range containers {
		err = p.docker.StopContainer(ctx, container.ID)
		if err != nil {
			return err
		}

		err = p.docker.RemoveContainer(ctx, container.ID)
		if err != nil {
			return err
		}
	}

	// удаление конфига роутера для проекта
	confFileName := utils.RouterConfFileName(name)

	fs := filesystem.NewFilesystem(utils.RouterConfDir)
	err = fs.FileRemove(confFileName)
	if err != nil {
		return err
	}

	router := services.NewRouter(p.docker)
	routerContainer, err := router.Find(ctx)
	if err != nil {
		return err
	}

	err = p.docker.RestartContainer(ctx, routerContainer.ID)
	if err != nil {
		return err
	}

	// удаление настроек для проекта
	projectSettings, err := settings.NewProjects()
	if err != nil {
		return err
	}

	projectSettings.RemoveProjectSettings(name)

	err = settings.Save(projectSettings)
	if err != nil {
		return err
	}

	return nil
}
