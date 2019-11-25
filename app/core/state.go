package core

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/project"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

type AppState struct {
	Router   *services.State
	Projects []*project.State
}

func (c *Core) State(ctx context.Context) (*AppState, error) {
	router := services.NewRouter(c.docker)

	rContainer, err := router.Find(ctx)
	rState := services.ContainerState(rContainer, utils.RouterName)

	prContainers, err := c.FindProjects(ctx)
	if err != nil {
		return nil, err
	}

	projects := make([]*project.State, 0, len(prContainers))
	for prName, prConts := range prContainers {
		if prName == utils.RouterName {
			continue
		}

		projects = append(projects, project.ProjectState(prName, prConts))
	}

	return &AppState{
		Router:   rState,
		Projects: projects,
	}, nil
}
