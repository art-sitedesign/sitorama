package core

import (
	"os"
	"text/template"

	"github.com/art-sitedesign/sitorama/app/utils"
)

func (c *Core) CreateSite(name string) error {
	// создание конфиг-фалйа для роутера
	err := createRouterConfig(name)
	if err != nil {
		return err
	}

	// создание контейнера

	return nil
}

func createRouterConfig(name string) error {
	tmpl := template.Must(template.ParseFiles(utils.RouterConfTemplate))

	err := os.MkdirAll(utils.RouterConfDir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(utils.RouterConfPath(name))
	if err != nil {
		return err
	}

	data := map[string]string{
		"Domain":    name,
		"Container": utils.ContainerName(name),
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
