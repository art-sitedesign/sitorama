package builder

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/art-sitedesign/sitorama/app/core/docker"
	"github.com/art-sitedesign/sitorama/app/core/services"
	"github.com/art-sitedesign/sitorama/app/utils"
)

const (
	mySQLDefaultPort = 3306
	mySQLDefaultUser = "root"
	mySQLForwardPort = "forwardPort"
	mySQLDataDir     = "mysql"
)

type MySQL struct {
	docker *docker.Docker
	name   string
	config Config
}

func NewMySQL(docker *docker.Docker, name string) *MySQL {
	return &MySQL{
		docker: docker,
		name:   name,
	}
}

func (m *MySQL) Name() string {
	return "MySQL"
}

func (m *MySQL) ConfigNames() []string {
	return []string{
		mySQLForwardPort,
	}
}

func (m *MySQL) ConfigByName(name string) (string, error) {
	switch name {
	case mySQLForwardPort:
		freePort := utils.FindNearPort(mySQLDefaultPort)
		return strconv.Itoa(freePort), nil
	default:
		return "", errors.New("unknown config " + name)
	}
}

func (m *MySQL) SetConfig(config Config) {
	m.config = config
}

func (m *MySQL) Checker() (string, error) {
	return renderChecker("mysql.php", map[string]string{
		"DBHost": m.alias(),
		"DBName": m.dbName(),
		"DBUser": m.user(),
		"DBPass": m.password(),
	})
}

func (m *MySQL) Info() map[string]string {
	return map[string]string{
		"Public host":  "127.0.0.1",
		"Public port":  m.portToForward(),
		"Private host": m.alias(),
		"Private port": services.MySQLPort,
		"User":         m.user(),
		"Password":     m.password(),
		"DB name":      m.dbName(),
	}
}

func (m *MySQL) Build(ctx context.Context) error {
	// поиск сети приложения
	network, err := m.docker.FindNetwork(ctx)
	if err != nil {
		return err
	}

	// сборка контейнера apache
	err = m.buildContainer(ctx, network.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) buildContainer(ctx context.Context, networkID string) error {
	dataPath, err := utils.ProjectVolumeFullPath(m.name, mySQLDataDir)
	if err != nil {
		return err
	}

	serviceMySQL := services.NewMySQL(
		m.docker,
		m.name,
		m.portToForward(),
		m.password(),
		m.dbName(),
		dataPath,
	)

	container, err := serviceMySQL.Find(ctx)
	if err != nil {
		return err
	}

	if container == nil {
		// если контейнер не найден - создаем
		cID, err := serviceMySQL.Create(ctx)
		if err != nil {
			return err
		}

		// подключаем к сети приложения созданный контейнер
		err = m.docker.ConnectNetwork(ctx, networkID, cID, []string{m.alias()})
		if err != nil {
			return err
		}

		// стартуем готовый контейнер
		err = m.docker.StartContainer(ctx, cID)
		if err != nil {
			return err
		}

		// установка аутентификации по дефолту с паролем (актуально для версий 8+)
		cmd := fmt.Sprintf(
			`mysql -u%s -P%s -p%s -e "alter user %s identified with mysql_native_password by '%s';"`,
			m.user(),
			services.MySQLPort,
			m.password(),
			m.user(),
			m.password(),
		)

		err = m.docker.ExecInContainer(ctx, cID, []string{cmd})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MySQL) alias() string {
	return fmt.Sprintf("%s.mysql", m.name)
}

func (m *MySQL) user() string {
	return mySQLDefaultUser
}

func (m *MySQL) password() string {
	return defaultPassword
}

func (m MySQL) dbName() string {
	return strings.Replace(m.name, ".", "_", -1)
}

func (m *MySQL) portToForward() string {
	confPort := m.config.String(mySQLForwardPort)
	if confPort != nil {
		return *confPort
	}

	return strconv.Itoa(mySQLDefaultPort)
}
