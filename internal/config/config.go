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
	JwtSecret string `yaml:"jwtSecret"`
	Port      int    `yaml:"port"`
}

type Config struct {
	UploadTemp  string `yaml:"uploadTemp"`
	UploadFinal string `yaml:"uploadFinal"`

	Serv     Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
}

func Load(filename string) (*Config, error) {
	// get jwtSecret & pgPassword from env var
	jwtSecret := os.Getenv("JWT_SECRET")
	pgPass := os.Getenv("POSTGRES_PASSWORD")

	if len(jwtSecret) == 0 {
		return nil, fmt.Errorf("jwt secret key not provided as env var")
	}

	if len(pgPass) == 0 {
		return nil, fmt.Errorf("postgres password not provided as env var")
	}

	rb, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(rb, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	config.Serv.JwtSecret = jwtSecret
	config.Postgres.Password = pgPass

	return &config, nil
}

func (c *Config) Port() string {
	return fmt.Sprintf(":%d", c.Serv.Port)
}

func (c *Config) JwtSecret() string {
	return c.Serv.JwtSecret
}
