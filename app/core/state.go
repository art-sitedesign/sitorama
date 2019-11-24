package core

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type AppState struct {
	Router   *services.State
	Services []*services.State
}

func (c *Core) State(ctx context.Context) (*AppState, error) {
	router := services.NewRouter(c.docker)

	rContainer, err := router.Find(ctx)
	rState := services.ContainerState(rContainer, utils.RouterName)

	projects, err := c.FindProjects(ctx)
	if err != nil {
		return nil, err
	}

	servs := make([]*services.State, 0, len(projects)*2)

	for projectName, containers := range projects {
		if projectName == utils.RouterName {
			continue
		}

		for _, container := range containers {
			servs = append(servs, services.ContainerState(&container, projectName))
		}
	}

	return &AppState{
		Router:   rState,
		Services: servs,
	}, nil
}
