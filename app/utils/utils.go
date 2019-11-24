package utils

import (
	"bytes"
	"fmt"
	"text/template"
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
