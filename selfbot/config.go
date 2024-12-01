package selfbot

import (
	"encoding/json"
	"errors"
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

func NewDefaultConfig(realToken string) Config {
	return Config{
		Token: realToken,
	}
}

func LoadConfig(jsonInBytes []byte) (Config, error) {
	config := NewDefaultConfig("")
	if err := json.Unmarshal(jsonInBytes, &config); err != nil {
		return Config{}, err
	}
	if config.Token == "" {
		return Config{}, errors.New("empty token specified")
	}
	return config, nil
}
