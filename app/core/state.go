package core

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type AppState struct {
	Router   *services.State
	Projects map[string][]*services.State
}

func (c *Core) State(ctx context.Context) (*AppState, error) {
	router := services.NewRouter(c.docker)

	rContainer, err := router.Find(ctx)
	rState := services.ContainerState(rContainer, utils.RouterName)

	projectContainers, err := c.FindProjects(ctx)
	if err != nil {
		return nil, err
	}

	projects := make(map[string][]*services.State)

	for projectName, cns := range projectContainers {
		if projectName == utils.RouterName {
			continue
		}

		for _, container := range cns {
			projects[projectName] = append(projects[projectName], services.ContainerState(&container, projectName))
		}
	}

	return &AppState{
		Router:   rState,
		Projects: projects,
	}, nil
}
