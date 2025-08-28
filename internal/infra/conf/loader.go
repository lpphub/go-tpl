package conf

import (
	"path/filepath"

	"github.com/lpphub/golib/env"
)

func LoadConfig() RConfig {
	var conf RConfig
	configFile := filepath.Join("config", "conf.yml")
	env.LoadConf(configFile, &conf)
	return conf
}
