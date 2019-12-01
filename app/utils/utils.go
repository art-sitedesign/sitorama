package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

// CreateRouterConfig создаст новый конфиг файл для роутера
func CreateRouterConfig(name string, containerAlias string) error {
	tmpl := template.Must(template.ParseFiles(RouterConfTemplate))

	err := os.MkdirAll(RouterConfDir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(RouterConfPath(name))
	if err != nil {
		return err
	}
	defer f.Close()

	data := map[string]string{
		"Domain":         name,
		"ContainerAlias": containerAlias,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

// AddHost добавит домен для локального хоста
func AddHost(name string) error {
	build := "linux"
	if isDarwin() {
		build = "macos"
	}

	cmd := exec.Command("make", "hm.add", "E="+build, "D="+name)
	return cmd.Run()
}

// RemoveHost удалит домен для локального хоста
func RemoveHost(name string) error {
	build := "linux"
	if isDarwin() {
		build = "macos"
	}

	cmd := exec.Command("make", "hm.rm", "E="+build, "D="+name)
	return cmd.Run()
}

// isDarwin вернет тип ОС. Костыль. Временно. Пока не хочу заморачиваться с нормальным сборщиком под ОС
func isDarwin() bool {
	return runtime.GOOS == "darwin"
}
