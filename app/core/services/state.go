package services

import "time"

type ServiceState struct {
	ID         string
	Name       string
	Active     bool
	CreateTime time.Time
}
