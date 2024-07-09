package config

import (
	"github.com/cristalhq/aconfig"
)

type Postgres struct {
	Host        string `env:"HOST"`
	Port        int    `env:"PORT"`
	Username    string `env:"USERNAME"`
	Password    string `env:"PASSWORD"`
	Database    string `env:"DATABASE"`
	SSLMode     string `env:"SSL_MODE" default:"disable"`
	SSLCertPath string `env:"SSL_CERT_PATH"`
}

type Config struct {
	Debug    bool     `env:"DEBUG"`
	Postgres Postgres `env:"POSTGRES"`
}

func MustLoad() *Config {
	cfg := Config{}

	err := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix: "UAUPDATE",
	}).Load()
	if err != nil {
		panic(err)
	}

	return &cfg
}
