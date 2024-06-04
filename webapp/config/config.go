package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config.yaml"`
	Server     struct {
		Host string `env:"HOST" env-default:"0.0.0.0"`
		Port int    `env:"PORT" env-default:"8080"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Insecure bool   `yaml:"insecure"`
	} `yaml:"database"`
	OAuth struct {
		ClientId     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
		RedirectUri  string `yaml:"redirectUri"`
	} `yaml:"oauth"`
	AuthStore struct {
		Name   string `yaml:"name"`
		Secret string `yaml:"secret"`
	} `yaml:"authStore"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
    Ollama struct {
        Url string `yaml:"url"`
        Model string `yaml:"model"`
    } `yaml:"ollama"`
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
