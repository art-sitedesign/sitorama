package core

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/services"
)

type AppState struct {
	RouterState *services.ServiceState
}

func (c *Core) State(ctx context.Context) *AppState {
	router := services.NewRouter(c.docker)

	return &AppState{
		RouterState: router.State(ctx),
	}
}
