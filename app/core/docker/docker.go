package docker

import "github.com/docker/docker/client"

type Docker struct {
	client *client.Client
}

func NewDocker() (*Docker, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Docker{client: c}, nil
}
