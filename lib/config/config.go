package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(target interface{}) error {
	path := flag.String("config", "config.yml", "The path to the configuration file")
	flag.Parse()

	f, err := os.Open(*path)

	if err != nil {
		return err
	}

	if err := yaml.NewDecoder(f).Decode(&target); err != nil {
		return err
	}

	return nil
}
