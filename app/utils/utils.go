package utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
)

func RouterConfPath(name string) string {
	return fmt.Sprintf("%s/%s.conf", RouterConfDir, name)
}

// ContainerName вернёт имя контейнера с суффиксом
func ContainerName(s string) string {
	return fmt.Sprintf("%s-%s", Prefix, s)
}

// RenderTemplateInBuffer рендерит шаблон в буфер
func RenderTemplateInBuffer(templatePath string, data interface{}) (*bytes.Buffer, error) {
	tmpl := template.Must(template.ParseFiles(templatePath))
	b := bytes.Buffer{}
	err := tmpl.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// ProjectNameFromContainer вернёт название проекта по названию контейнера
func ProjectNameFromContainer(c *types.Container) string {
	//todo: сделать нормально!
	cName := c.Names[0][1:]
	cName = strings.Replace(cName, Prefix+"-", "", -1)
	sp := strings.Split(cName, "_")

	if len(sp) > 1 {
		sp = sp[:len(sp)-1]
	}

	return strings.Join(sp, "_")
}
