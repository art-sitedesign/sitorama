package settings

const projectsSettingsFileName = "projects-config.json"

// так себе решение :) но нет времени сделать иначе
type projectsSettings map[string]map[string]map[string]string

type Projects struct {
	Projects projectsSettings
}

func NewProjects() (*Projects, error) {
	p := &Projects{}
	err := load(p)
	if err != nil {
		return nil, err
	}

	if p.Projects == nil {
		p.Projects = make(projectsSettings)
	}

	return p, nil
}

func (p *Projects) FileName() string {
	return projectsSettingsFileName
}

func (p *Projects) AddBuilderSettings(project string, builder string, settings map[string]string) {
	pr, ok := p.Projects[project]
	if !ok {
		pr = make(map[string]map[string]string)
	}

	pr[builder] = settings

	p.Projects[project] = pr
}

func (p *Projects) RemoveProjectSettings(project string) {
	delete(p.Projects, project)
}
