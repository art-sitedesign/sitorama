package project

import (
	"context"
	"os"

	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

func (p *Project) Remove(ctx context.Context, name string) error {
	containers, err := p.docker.FindContainers(ctx, name)
	if err != nil {
		return err
	}

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

	routerConfPath := utils.RouterConfPath(name)
	err = os.Remove(routerConfPath)
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

	return nil
}
