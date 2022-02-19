package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Builders []struct {
		Name     string   `toml:"name"`
		Url      string   `toml:"url"`
		Commands []string `toml:"commands"`
	} `toml:"builders"`
}

func GetConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
