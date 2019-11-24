package core

import "context"

func (c *Core) ContainerRestart(ctx context.Context, cID string) error {
	return c.docker.RestartContainer(ctx, cID)
}

func (c *Core) ContainerStop(ctx context.Context, cID string) error {
	return c.docker.StopContainer(ctx, cID)
}

func (c *Core) ContainerStart(ctx context.Context, cID string) error {
	return c.docker.StartContainer(ctx, cID)
}
