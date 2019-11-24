package core

import (
	"context"
	"fmt"

	"github.com/art-sitedesign/sitorama/app/core/services"
)

type AppState struct {
	Router   *services.State
	Services []*services.State
}

func (c *Core) State(ctx context.Context) *AppState {
	router := services.NewRouter(c.docker)

	//todo: найти контейнеры, понять сайты и распарсить
	snginx := services.NewSiteNginx(c.docker, "test.loc", fmt.Sprintf("%s.phpfpm", "test.loc"))
	sphpfpm := services.NewSitePHPFPM(c.docker, "test.loc")

	servs := make([]*services.State, 0, 2)
	servs = append(servs, services.ServiceState(ctx, snginx), services.ServiceState(ctx, sphpfpm))

	return &AppState{
		Router:   services.ServiceState(ctx, router),
		Services: servs,
	}
}
