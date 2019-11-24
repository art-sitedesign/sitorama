package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// FindImage найдет имайдж
func (d *Docker) FindImage(ctx context.Context, name string) ([]types.ImageSummary, error) {
	f := filters.NewArgs()
	f.Add("reference", name)
	options := types.ImageListOptions{
		Filters: f,
	}

	return d.client.ImageList(ctx, options)
}

// PullImage скачает имайдж локально
func (d *Docker) PullImage(ctx context.Context, name string) error {
	//todo: в либе могут вмержить фикс и тогда канонический урл не нужен будет
	canonicalName := fmt.Sprintf("docker.io/%s", name)

	options := types.ImagePullOptions{
		All:           false,
		RegistryAuth:  "",
		PrivilegeFunc: nil,
	}
	rc, err := d.client.ImagePull(ctx, canonicalName, options)
	if err != nil {
		return err
	}

	res, err := d.client.ImageLoad(ctx, rc, false)

	_ = res

	_ = rc.Close()

	return nil
}
