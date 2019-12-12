package project

import (
	"context"

	"github.com/art-sitedesign/sitorama/app/core/builder"
	"github.com/art-sitedesign/sitorama/app/core/settings"
)

func (p *Project) Create(ctx context.Context, builders []builder.Builder) error {
	projectsSettings, err := settings.NewProjects()
	if err != nil {
		return err
	}

	for _, b := range builders {
		err := b.Build(ctx)
		if err != nil {
			return err
		}

		projectsSettings.AddBuilderSettings(p.name, b.Name(), b.Info())
	}

	err = settings.Save(projectsSettings)
	if err != nil {
		return err
	}

	return nil
}
