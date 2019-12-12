package project

import (
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/core/settings"
	"github.com/docker/docker/api/types"
)

type State struct {
	Name     string
	Services []*services.State
	Info     map[string]map[string]string
}

func NewState(n string, s []*services.State, i map[string]map[string]string) *State {
	return &State{
		Name:     n,
		Services: s,
		Info:     i,
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

func ProjectState(name string, containers []types.Container, projectsSettings *settings.Projects) *State {
	states := make([]*services.State, 0, len(containers))
	for _, container := range containers {
		states = append(states, services.ContainerState(&container, name))
	}

	return NewState(name, states, projectsSettings.Projects[name])
}
