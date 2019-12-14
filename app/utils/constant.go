package utils

const (
	Prefix = "sitorama"

	RouterName         = "router"
	RouterConfDir      = "volumes/router/nginx"
	RouterConfTemplate = "app/templates/nginx/router.conf"

	SiteNginxBaseTemplate   = "app/templates/nginx/site-nginx.conf"
	SiteNginxServerTemplate = "app/templates/nginx/site-server.conf"
	ApacheServerTemplate    = "app/templates/apache/server.conf"
	PHPIniTemplate          = "app/templates/php/php.ini"
	PHPLibsTemplate         = "app/templates/php/libs.ini"

	CheckersDir      = "app/templates/php/checkers"
	IndexPHPTemplate = "app/templates/php/index.php"

	ProjectVolumesPath = "sitorama_volumes"

	linuxBinName = "linux-amd64"
	macOSBinName = "darwin-amd64"
)
