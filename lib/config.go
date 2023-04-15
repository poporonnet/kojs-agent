package lib

import (
	"github.com/go-yaml/yaml"
	"github.com/mct-joken/jkojs-agent/types"
	"os"
)

var Config types.Config

func LoadConfig() error {
	f, err := os.ReadFile("config.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, &Config)
	if err != nil {
		return err
	}
	return nil
}
