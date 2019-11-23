package core

import "pet-projects/sitorama/app/core/docker"

type Core struct {
	docker *docker.Docker
}

func NewCore() (*Core, error) {
	d, err := docker.NewDocker()
	if err != nil {
		return nil, err
	}

	return &Core{docker: d}, nil
}
