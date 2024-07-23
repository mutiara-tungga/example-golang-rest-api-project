package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type envConfig struct {
	AppPort string `env:"APP_PORT"`

	DatabaseHost        string `env:"DATABASE_HOST"`
	DatabasePort        string `env:"DATABASE_PORT"`
	DatabaseUser        string `env:"DATABASE_USER"`
	DatabasePass        string `env:"DATABASE_PASS"`
	DatabaseName        string `env:"DATABASE_NAME"`
	DatabaseMaxOpenConn int    `env:"DATABASE_MAX_OPEN_CONN"`
	DatabaseMaxIdleConn int    `env:"DATABASE_MAX_IDLE_CONN"`
	DatabaseMaxLifeTime int    `env:"DATABASE_MAX_LIFE_TIME"`

	Version string `env:"VERSION"`
}

var (
	envCfg *envConfig
	once   sync.Once
)

func Get() *envConfig {
	return envCfg
}

func LoadEnvConfig() {
	once.Do(func() {
		_ = godotenv.Overload()
		c := new(envConfig)
		if err := env.Parse(c); err != nil {
			panic(err)
		}
		envCfg = c
		// cfg.Version = buildutil.Version // TODO: get when build
	})
}
