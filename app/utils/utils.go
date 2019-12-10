package utils

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"

	"github.com/art-sitedesign/sitorama/app/core/filesystem"
	"github.com/art-sitedesign/sitorama/app/core/settings"
)

func RouterConfFileName(name string) string {
	return fmt.Sprintf("%s.conf", name)
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

	confFileName := RouterConfFileName(name)

	fs := filesystem.NewFilesystem(RouterConfDir)
	f, err := fs.FileOpenForce(confFileName, os.O_RDWR)
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

// ProjectFullPath вернёт абсолютный путь до корневой директории проекта
func ProjectFullPath(name string) (string, error) {
	s, err := settings.NewApp()
	if err != nil {
		return "", err
	}

	fs := filesystem.NewFilesystem(s.ProjectsRoot)
	fs.AddDir(name)

	return fs.FullPath()
}

// ProjectVolumesFullPath вернёт абсолютный путь до директории с вольюмами проекта
func ProjectVolumeFullPath(projectName string, volumeName string) (string, error) {
	s, err := settings.NewApp()
	if err != nil {
		return "", err
	}

	fs := filesystem.NewFilesystem(s.ProjectsRoot)
	fs.AddDir(projectName)
	fs.AddDir(ProjectVolumesPath)
	fs.AddDir(volumeName)

	if !fs.Exist() {
		err = fs.Create()
		if err != nil {
			return "", err
		}
	}

	return fs.FullPath()
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

// FindNearPort найдет ближайший свободный порт к переданному увеличивая значение при каждой попытке
func FindNearPort(from int) int {
	for !IsPortFree(from) {
		from++
	}

	return from
}

// IsPortFree проверить свободен ли порт
func IsPortFree(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return false
	}

	_ = ln.Close()

	return true
}

// isDarwin вернет тип ОС. Костыль. Временно. Пока не хочу заморачиваться с нормальным сборщиком под ОС
func isDarwin() bool {
	return runtime.GOOS == "darwin"
}
