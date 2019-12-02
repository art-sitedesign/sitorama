package settings

import (
	"encoding/json"
)

const appSettingsFileName = "app-config.json"

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

	err = fileSystem().FileWrite(appSettingsFileName, data)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) load() error {
	data, err := fileSystem().FileRead(appSettingsFileName)
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
