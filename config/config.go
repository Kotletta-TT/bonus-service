package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

const (
	defaultLocalAddr     = "localhost:8080"
	defaultLocalDB       = "postgres://bonus:bonus@localhost:5432/bonusdb"
	defaultLocalLogLevel = "DEBUG"
	defaultLocalSecret   = "79675945011f9de2372d2ddd9510907938ed586896e680a67a36bb1e86e54839"
)

type Config struct {
	ServAddr   string `env:"RUN_ADDRESS"`
	DSN        string `env:"DATABASE_URI"`
	AccuralURL string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel   string `env:"LOG_LEVEL"`
	SecretKey  string `env:"SECRET_KEY"`
}

func New() *Config {
	config := Config{}
	flag.StringVar(&config.ServAddr, "a", defaultLocalAddr, "Address:port server")
	flag.StringVar(&config.DSN, "d", defaultLocalDB, "DB URL example: postgres://bonus:bonus@localhost:5432/bonusdb")
	flag.StringVar(&config.AccuralURL, "r", "", "URL accural system")
	flag.StringVar(&config.LogLevel, "level", defaultLocalLogLevel, "Log level")
	flag.StringVar(&config.SecretKey, "secret", defaultLocalSecret, "Secret key for generate JWT token")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
