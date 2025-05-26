package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func prepare(configPath string) error {
	sc := ServerConfig{
		Port: 8000,
	}

	conf := Config{
		Server: sc,
	}

	wb, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to marshal config into yaml: %v", err)
	}

	if err := os.WriteFile(configPath, wb, 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}

func Load(configPath string) (*Config, error) {
	rb, err := os.ReadFile(configPath)
	if err != nil {
		if perr := prepare(configPath); perr != nil {
			return nil, perr
		}

		// try to read again
		rb, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", configPath, err)
		}
	}

	var conf Config
	if err := yaml.Unmarshal(rb, &conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return &conf, nil
}

func (c *Config) ServerPort() string {
	return fmt.Sprintf("0.0.0.0:%d", c.Server.Port)
}
