package core

import (
	"context"

	"github.com/docker/docker/api/types"

	"github.com/art-sitedesign/sitorama/app/core/builder"
	"github.com/art-sitedesign/sitorama/app/core/project"
	"github.com/art-sitedesign/sitorama/app/utils"
)

// FindProjects найдет контейнеры проектов в вернёт их сгруппированными по проекту
func (c *Core) FindProjects(ctx context.Context) (map[string][]types.Container, error) {
	containers, err := c.docker.FindContainers(ctx, "")
	if err != nil {
		return nil, err
	}

	result := make(map[string][]types.Container)
	for _, container := range containers {
		projectName := utils.ProjectNameFromContainer(&container)
		result[projectName] = append(result[projectName], container)
	}

	return result, nil
}

// CreateProject создаст проект
func (c *Core) CreateProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	//todo: понять какие билдеры нужны...
	builders := make([]builder.Builder, 0, 1)

	nginxPHPFPM := builder.NewNginxPHPFPM(c.docker, name)
	builders = append(builders, nginxPHPFPM)

	err := pr.Create(ctx, builders)
	if err != nil {
		return err
	}

	err = utils.AddHost(name)
	if err != nil {
		return ErrorCantChangeHosts
	}

	return nil
}

// StartProject запустит контейнеры проекта
func (c *Core) StartProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	return pr.Start(ctx, name)
}

// StopProject остановит контейнеры проекта
func (c *Core) StopProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	return pr.Stop(ctx, name)
}

// RemoveProject удалит проект
func (c *Core) RemoveProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	err := pr.Remove(ctx, name)
	if err != nil {
		return err
	}

	err = utils.RemoveHost(name)
	if err != nil {
		return ErrorCantChangeHosts
	}

	return nil
}
