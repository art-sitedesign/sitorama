package project

import "context"

func (p *Project) Start(ctx context.Context, name string) error {
	containers, err := p.docker.FindContainers(ctx, name)
	if err != nil {
		return err
	}

	for _, container := range containers {
		err = p.docker.StartContainer(ctx, container.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
