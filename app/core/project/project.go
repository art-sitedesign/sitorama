package project

import "github.com/art-sitedesign/sitorama/app/core/docker"

type Project struct {
	docker *docker.Docker
}

func NewProject(d *docker.Docker) *Project {
	return &Project{
		docker: d,
	}
}
