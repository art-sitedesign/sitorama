package utils

import (
	"fmt"
)

func RouterConfPath(name string) string {
	return fmt.Sprintf("%s/%s.conf", RouterConfDir, name)
}

// ContainerName вернёт имя контейнера с суффиксом
func ContainerName(s string) string {
	return fmt.Sprintf("%s-%s", Prefix, s)
}
