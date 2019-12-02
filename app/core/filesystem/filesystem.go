package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Filesystem struct {
	path string
}

func NewFilesystem(path string) *Filesystem {
	return &Filesystem{
		path: strings.Trim(path, "/"),
	}
}

func (fs *Filesystem) Exist() bool {
	return isExist(fs.path)
}

func (fs *Filesystem) FileExist(fileName string) bool {
	return isExist(fs.filePath(fileName))
}

func (fs *Filesystem) Create() error {
	return os.MkdirAll(fs.path, 0755)
}

func (fs *Filesystem) FileCreate(fileName string) (*os.File, error) {
	if !fs.Exist() {
		err := fs.Create()
		if err != nil {
			return nil, err
		}
	}

	return os.Create(fs.filePath(fileName))
}

func (fs *Filesystem) FullPath() (string, error) {
	return fullPath(fs.path)
}

func (fs *Filesystem) FileFullPath(fileName string) (string, error) {
	return fullPath(fs.filePath(fileName))
}

func (fs *Filesystem) FileOpenForce(fileName string, flag int) (*os.File, error) {
	fp := fs.filePath(fileName)

	f, err := os.OpenFile(fp, flag, 0755)
	if os.IsNotExist(err) {
		f, err = os.Create(fp)
	}

	return f, err
}

func (fs *Filesystem) FileWrite(fileName string, data []byte) error {
	return ioutil.WriteFile(fs.filePath(fileName), data, 0755)
}

func (fs *Filesystem) FileRead(fileName string) ([]byte, error) {
	f, err := fs.FileOpenForce(fileName, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}

func (fs *Filesystem) FileRemove(fileName string) error {
	return os.Remove(fs.filePath(fileName))
}

func (fs *Filesystem) filePath(fileName string) string {
	return fmt.Sprintf("%s/%s", fs.path, strings.TrimLeft(fileName, "/"))
}

func isExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func fullPath(name string) (string, error) {
	return filepath.Abs(name)
}
