package settings

import (
	"encoding/json"

	"github.com/art-sitedesign/sitorama/app/core/filesystem"
)

const (
	settingsPath = "app"
)

type Settings interface {
	FileName() string
}

func Save(s Settings) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = fileSystem().FileWrite(s.FileName(), data)
	if err != nil {
		return err
	}

	return nil
}

func fileSystem() *filesystem.Filesystem {
	return filesystem.NewFilesystem(settingsPath)
}

func load(s Settings) error {
	data, err := fileSystem().FileRead(s.FileName())
	if err != nil {
		return err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, s)
		if err != nil {
			return err
		}
	}

	return nil
}
