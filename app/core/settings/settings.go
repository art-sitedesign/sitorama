package settings

import "github.com/art-sitedesign/sitorama/app/core/filesystem"

const (
	settingsPath = "app"
)

func fileSystem() *filesystem.Filesystem {
	return filesystem.NewFilesystem(settingsPath)
}
