package project

import "context"

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

	return nil
}
