package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

const (
	NginxPHPFPMConfigNginx = "nginx-config"
)

type NginxPHPFPM struct {
	docker *docker.Docker
	name   string
	config Config
}

func NewNginxPHPFPM(d *docker.Docker, n string) Builder {
	return &NginxPHPFPM{
		docker: d,
		name:   n,
	}
}

func (np *NginxPHPFPM) Name() string {
	return "Nginx+PHP-FPM"
}

func (np *NginxPHPFPM) ConfigNames() []string {
	return []string{
		NginxPHPFPMConfigNginx,
	}
}

func (np *NginxPHPFPM) ConfigByName(name string) (string, error) {
	switch name {
	case NginxPHPFPMConfigNginx:
		nConf, err := np.nginxConfig()
		if err != nil {
			return "", err
		}
		return nConf, nil
	default:
		return "", errors.New("unknown config " + name)
	}
}

func (np *NginxPHPFPM) PrepareConfig() (Config, error) {
	conf := Config{}

	for _, name := range np.ConfigNames() {
		c, err := np.ConfigByName(name)
		if err != nil {
			return nil, err
		}
		conf[name] = c
	}

	return conf, nil
}

func (np *NginxPHPFPM) SetConfig(config Config) {
	np.config = config
}

func (np *NginxPHPFPM) Build(ctx context.Context) error {
	ngAlias, _ := np.aliases()

	// создание конфиг-фалйа для роутера
	err := utils.CreateRouterConfig(np.name, ngAlias)
	if err != nil {
		return err
	}

	// поиск сети приложения
	network, err := np.docker.FindNetwork(ctx)
	if err != nil {
		return err
	}

	// сборка контейнера PHP-FPM
	err = np.buildPHPFPM(ctx, network.ID)
	if err != nil {
		return err
	}

	err = np.buildNginx(ctx, network.ID)
	if err != nil {
		return err
	}

	// перезапуск роутера для применения конфига
	r := services.NewRouter(np.docker)
	router, err := r.Find(ctx)
	if err != nil {
		return err
	}

	err = np.docker.RestartContainer(ctx, router.ID)
	if err != nil {
		return err
	}

	return nil
}

func (np *NginxPHPFPM) buildPHPFPM(ctx context.Context, networkID string) error {
	_, pfAlias := np.aliases()
	// поиск PHP-FPM контейнера
	sitePHPFPM := services.NewSitePHPFPM(np.docker, np.name)
	pfContainer, err := sitePHPFPM.Find(ctx)
	if err != nil {
		return err
	}

	if pfContainer == nil {
		// если контейнер не найден - создаем
		pfCID, err := sitePHPFPM.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = np.docker.ConnectNetwork(ctx, networkID, pfCID, []string{pfAlias})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = np.docker.StartContainer(ctx, pfCID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (np *NginxPHPFPM) buildNginx(ctx context.Context, networkID string) error {
	ngAlias, pfAlias := np.aliases()

	var nConfP *string
	nConf, ok := np.config[NginxPHPFPMConfigNginx]
	if ok {
		nConfP = &nConf
	} else {
		nConfP = nil
	}

	siteNginx := services.NewSiteNginx(np.docker, np.name, pfAlias, nConfP)
	container, err := siteNginx.Find(ctx)
	if err != nil {
		return err
	}

	if container == nil {
		// если контейнер не найден - создаем
		cID, err := siteNginx.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = np.docker.ConnectNetwork(ctx, networkID, cID, []string{ngAlias})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = np.docker.StartContainer(ctx, cID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (np *NginxPHPFPM) nginxConfig() (string, error) {
	_, pfAlias := np.aliases()
	siteNginx := services.NewSiteNginx(np.docker, np.name, pfAlias, nil)
	nConf, err := siteNginx.RenderConfig()
	if err != nil {
		return "", err
	}

	return nConf.String(), nil
}

func (np *NginxPHPFPM) aliases() (string, string) {
	pfAlias := fmt.Sprintf("%s.phpfpm", np.name)
	ngAlias := fmt.Sprintf("%s.nginx", np.name)

	return ngAlias, pfAlias
}
