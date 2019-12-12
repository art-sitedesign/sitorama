package settings

const appSettingsFileName = "app-config.json"

type App struct {
	ProjectsRoot string
}

func NewApp() (*App, error) {
	a := &App{}
	err := load(a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) FileName() string {
	return appSettingsFileName
}
