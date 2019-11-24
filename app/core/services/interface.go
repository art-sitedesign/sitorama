package services

import (
	"context"

	"github.com/docker/docker/api/types"
)

type Service interface {
	Find(ctx context.Context) (*types.Container, error)
	Create(ctx context.Context) (string, error)
	ContainerName() string
	State(ctx context.Context) *ServiceState
}
