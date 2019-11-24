package core

import (
	"context"

	"github.com/docker/docker/api/types"

	"github.com/art-sitedesign/sitorama/app/utils"
)

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
