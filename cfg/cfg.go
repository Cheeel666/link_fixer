package cfg

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func NewCfg(path string) (*Config, error) {
	var cfg *Config
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
