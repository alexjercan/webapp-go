package webapp

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config.yaml"`
	Server     struct {
        Host string `env:"HOST" env-default:""`
		Port int `env:"PORT" env-default:"8080"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Insecure bool   `yaml:"insecure"`
	} `yaml:"database"`
}

func LoadConfig() (cfg Config, err error) {
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return
	}

	err = cleanenv.ReadConfig(cfg.ConfigPath, &cfg)
	if err != nil {
		return
	}

	return
}
