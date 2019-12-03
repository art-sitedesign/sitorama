package builder

import (
	"context"
)

type Builder interface {
	Name() string
	ConfigNames() []string
	ConfigByName(name string) (string, error)
	SetConfig(config Config)
	Build(ctx context.Context) error
}

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

type Config map[string]string

func (c Config) String(key string) *string {
	var res *string
	conf, ok := c[key]
	if ok {
		res = &conf
	} else {
		res = nil
	}

	return res
}

func PrepareConfig(builder Builder) (Config, error) {
	conf := Config{}

	for _, name := range builder.ConfigNames() {
		c, err := builder.ConfigByName(name)
		if err != nil {
			return nil, err
		}
		conf[name] = c
	}

	return conf, nil
}
