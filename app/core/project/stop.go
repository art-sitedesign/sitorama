package project

import "context"

func (p *Project) Stop(ctx context.Context, name string) error {
	containers, err := p.docker.FindContainers(ctx, name)
	if err != nil {
		return err
	}

	for _, container := range containers {
		err = p.docker.StopContainer(ctx, container.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
