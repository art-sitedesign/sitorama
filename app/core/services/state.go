package services

import (
	"context"
	"time"
)

type State struct {
	ID         string
	Name       string
	Active     bool
	CreateTime time.Time
}

// ServiceState вернет состояние сервиса
func ServiceState(ctx context.Context, service Service) *State {
	container, err := service.Find(ctx)

	ss := &State{}

	if err == nil && container != nil {
		ss.ID = container.ID[:12]
		ss.Name = container.Names[0][1:]
		ss.Active = container.State == "running"
		ss.CreateTime = time.Unix(container.Created, 0)
	}

	return ss
}
