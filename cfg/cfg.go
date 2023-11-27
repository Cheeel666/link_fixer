package cfg

import (
	"errors"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func NewCfg() (cfg Config, err error) {
	cfg.Token = os.Getenv("TG_TOKEN")
	if cfg.Token == "" {
		err = errors.New("missing telegram token")
		return
	}
	return
}
