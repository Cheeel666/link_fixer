package cfg

import (
	"errors"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func NewCfg() (*Config, error) {
	cfg := &Config{
		Token: os.Getenv("TG_TOKEN"),
	}

	if cfg.Token == "" {
		return nil, errors.New("missing telegram token")
	}

	return cfg, nil
}
