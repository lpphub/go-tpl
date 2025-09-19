package conf

type RConfig struct {
	Mysql mysqlConf `yaml:"mysql"`
	Redis redisConf `yaml:"redis"`
	Log   logConf   `yaml:"log"`
}

type mysqlConf struct {
	Addr     string `yaml:"addr"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type redisConf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type logConf struct {
	Path string `yaml:"path"`
}
