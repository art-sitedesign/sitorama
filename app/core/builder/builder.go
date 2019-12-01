package builder

import (
	"context"
)

type Builder interface {
	Name() string
	ConfigNames() []string
	ConfigByName(name string) (string, error)
	PrepareConfig() (Config, error)
	SetConfig(config Config)
	Build(ctx context.Context) error
}

type Config map[string]string

const (
	BuilderNginxPHPFPM = 1
	BuilderApache      = 2
	BuilderPostgres    = 3
	BuilderMySQL       = 4
	BuilderRedis       = 5
	BuilderMemcached   = 6
)

var (
	WebserverBuilders = map[int]string{
		BuilderNginxPHPFPM: "Nginx+PHP-FPM",
		BuilderApache:      "Apache",
	}

	DatabaseBuilders = map[int]string{
		BuilderPostgres: "Postgres",
		BuilderMySQL:    "MySQL",
	}

	CacheBuilders = map[int]string{
		BuilderRedis:     "Redis",
		BuilderMemcached: "Memcached",
	}
)
