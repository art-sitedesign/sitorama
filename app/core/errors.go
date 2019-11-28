package core

import "errors"

var (
	ErrorCantChangeHosts = errors.New("не удалось изменить файл /etc/hosts")
)
