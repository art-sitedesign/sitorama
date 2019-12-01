package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const appSettingsPath = "app/app-config.json"

type App struct {
	ProjectsRoot string
}

func NewApp() (*App, error) {
	a := &App{}
	err := a.load()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Save() error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(appSettingsPath, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) load() error {
	f, err := a.file()
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) file() (*os.File, error) {
	f, err := os.OpenFile(appSettingsPath, os.O_RDWR, 0755)
	if os.IsNotExist(err) {
		f, err = os.Create(appSettingsPath)
	}

	return f, err
}
