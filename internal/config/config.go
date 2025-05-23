package config

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

type Postgres struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databaseName"`
	SslMode      string `yaml:"sslmode"`
}

type Server struct {
	Port int `yaml:"Port"`
}

type Config struct {
	UploadTemp  string `yaml:"uploadTemp"`
	UploadFinal string `yaml:"uploadFinal"`

	Serv     Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
}

func Load(filename string) (*Config, error) {
	rb, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(rb, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return &config, nil
}

func (c *Config) Port() string {
	return fmt.Sprintf(":%d", c.Serv.Port)
}
