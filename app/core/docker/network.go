package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

// CreateNetwork создаст новую подсеть для проекта
func (d *Docker) CreateNetwork(ctx context.Context) (string, error) {
	options := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		EnableIPv6:     false,
		IPAM:           nil,
		Internal:       false,
		Attachable:     false,
		Options:        nil,
		Labels:         nil,
	}

	res, err := d.client.NetworkCreate(ctx, prefix, options)
	if err != nil {
		return "", err
	}

	//todo: log warnings
	fmt.Println(res.Warning)

	return res.ID, nil
}

// FindNetwork найдет подсеть проекта
func (d *Docker) FindNetwork(ctx context.Context) (*types.NetworkResource, error) {
	args := filters.NewArgs()
	args.Add("name", prefix)
	options := types.NetworkListOptions{Filters: args}

	networks, err := d.client.NetworkList(ctx, options)
	if err != nil {
		return nil, err
	}

	if len(networks) == 0 {
		return nil, nil
	}

	return &networks[0], nil
}

// ConnectNetwork подключит контейнер к подсети
func (d *Docker) ConnectNetwork(ctx context.Context, networkID string, containerID string, aliases []string) error {
	res, err := d.client.NetworkInspect(ctx, networkID)
	if err != nil {
		return err
	}

	if _, ok := res.Containers[containerID]; ok {
		// контейнер уже в сети
		return nil
	}

	config := &network.EndpointSettings{
		IPAMConfig:          nil,
		Links:               nil,
		Aliases:             aliases,
		NetworkID:           "",
		EndpointID:          "",
		Gateway:             "",
		IPAddress:           "",
		IPPrefixLen:         0,
		IPv6Gateway:         "",
		GlobalIPv6Address:   "",
		GlobalIPv6PrefixLen: 0,
		MacAddress:          "",
	}

	return d.client.NetworkConnect(ctx, networkID, containerID, config)
}
