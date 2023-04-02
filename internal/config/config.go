package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Grpc struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Mysqldatabase struct {
		User     string `yaml:"user"`
		Host     string `yaml:"dbhost"`
		Port     string `yaml:"dbport"`
		Password string `yaml:"password"`
	}

	Crypt struct {
		Secret string `yaml:"secret"`
	}

	Environment string
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	err = d.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
