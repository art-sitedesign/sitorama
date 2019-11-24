package services

import (
	"time"

	"github.com/docker/docker/api/types"
)

type State struct {
	ID            string
	ServiceName   string
	ContainerName string
	Active        bool
	CreateTime    time.Time
}

func (s *State) CanStart() bool {
	return s.Active == false && s.ID != ""
}

// ContainerState вернёт состояние контейнера
func ContainerState(container *types.Container, projectName string) *State {
	ss := &State{}

	if container != nil {
		ss.ID = container.ID[:12]
		ss.ServiceName = projectName
		ss.ContainerName = container.Names[0][1:]
		ss.Active = container.State == "running"
		ss.CreateTime = time.Unix(container.Created, 0)
	}

	return ss
}
