package project

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/builder"
)

func (p *Project) Create(ctx context.Context, builders []builder.Builder) error {
	for _, bldr := range builders {
		err := bldr.Build(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
