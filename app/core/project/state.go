package project

import (
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/docker/docker/api/types"
)

type State struct {
	Name     string
	Services []*services.State
}

func NewState(n string, s []*services.State) *State {
	return &State{
		Name:     n,
		Services: s,
	}
}

func (s *State) Active() bool {
	a := true
	for _, ss := range s.Services {
		if !ss.Active {
			a = false
			break
		}
	}

	return a
}

func ProjectState(name string, containers []types.Container) *State {
	states := make([]*services.State, 0, len(containers))
	for _, container := range containers {
		states = append(states, services.ContainerState(&container, name))
	}

	return NewState(name, states)
}
