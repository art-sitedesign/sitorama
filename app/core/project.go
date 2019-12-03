package core

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"

	"github.com/art-sitedesign/sitorama/app/core/builder"
	"github.com/art-sitedesign/sitorama/app/core/filesystem"
	"github.com/art-sitedesign/sitorama/app/core/project"
	"github.com/art-sitedesign/sitorama/app/core/settings"
	"github.com/art-sitedesign/sitorama/app/models"
	"github.com/art-sitedesign/sitorama/app/utils"
)

// FindProjects найдет контейнеры проектов в вернёт их сгруппированными по проекту
func (c *Core) FindProjects(ctx context.Context) (map[string][]types.Container, error) {
	containers, err := c.docker.FindContainers(ctx, "")
	if err != nil {
		return nil, err
	}

	result := make(map[string][]types.Container)
	for _, container := range containers {
		projectName := utils.ProjectNameFromContainer(&container)
		result[projectName] = append(result[projectName], container)
	}

	return result, nil
}

// CreateProject создаст проект
func (c *Core) CreateProject(ctx context.Context, model *models.ProjectCreate, builders []builder.Builder) error {
	indexFileName := "index.php"

	appSettings, err := settings.NewApp()
	if err != nil {
		return err
	}

	if appSettings.ProjectsRoot == "" {
		return errors.New("projects root is not set")
	}

	fs := filesystem.NewFilesystem(appSettings.ProjectsRoot)

	// создаем корень для проектов
	err = fs.Create()
	if err != nil {
		return err
	}

	// создаем директорию для проекта
	fs.AddDir(model.Domain)
	err = fs.Create()
	if err != nil {
		return err
	}

	// создаем директорию с точкой входа в проект
	fs.AddDir(model.EntryPoint)
	err = fs.Create()
	if err != nil {
		return err
	}

	// создаем шаблонный index.php
	f, err := fs.FileCreate(indexFileName)
	if err != nil {
		return err
	}

	if f != nil {
		// если файла небыло и он только что создался - запишем в него шаблон
		data := map[string]string{"Name": model.Domain}
		b, err := utils.RenderTemplateInBuffer(utils.IndexPHPTemplate, data)
		if err != nil {
			return err
		}

		err = fs.FileWrite(indexFileName, b.Bytes())
		if err != nil {
			return err
		}
	}

	// создаем контейнеры проекта
	pr := project.NewProject(c.docker)

	err = pr.Create(ctx, builders)
	if err != nil {
		return err
	}

	// добавляем домен в /etc/hosts
	err = utils.AddHost(model.Domain)
	if err != nil {
		return ErrorCantChangeHosts
	}

	return nil
}

// StartProject запустит контейнеры проекта
func (c *Core) StartProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	return pr.Start(ctx, name)
}

// StopProject остановит контейнеры проекта
func (c *Core) StopProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	return pr.Stop(ctx, name)
}

// RemoveProject удалит проект
func (c *Core) RemoveProject(ctx context.Context, name string) error {
	pr := project.NewProject(c.docker)

	err := pr.Remove(ctx, name)
	if err != nil {
		return err
	}

	err = utils.RemoveHost(name)
	if err != nil {
		return ErrorCantChangeHosts
	}

	return nil
}
