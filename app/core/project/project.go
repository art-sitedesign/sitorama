package project

import "github.com/art-sitedesign/sitorama/app/core/docker"

type Project struct {
	docker *docker.Docker
	name   string
}

func NewProject(d *docker.Docker, name string) *Project {
	return &Project{
		docker: d,
		name:   name,
	}
}
