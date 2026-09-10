package config

import (
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		Base Base
		DB   DB
	}

	Base struct {
		Address       string `env:"ADRESS,required"`
		DataDir       string `env:"DATADIR,required"`
		AdminToken    string `env:"ADMIN_TOKEN,required"`
		CacheTTL      int64  `env:"CACHE_TTL,required"`
		MaxUploadSize int64  `env:"MAX_UPLOAD_SIZE,required"`
		MaxCacheItem  int64  `env:"MAX_CACHE_ITEM,required"`
	}
	DB struct {
		Host     string `env:"POSTGRES_HOST,required"`
		Port     string `env:"POSTGRES_PORT,required"`
		DbName   string `env:"POSTGRES_DB,required"`
		User     string `env:"POSTGRES_USER,required"`
		Password string `env:"POSTGRES_PASSWORD,required"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
