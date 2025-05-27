package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type PostgresConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databaseName"`
	SslMode      string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port       int    `yaml:"port"`
	JwtSecret  string `yaml:"jwtSecret"`
	UploadPath string `yaml:"uploadPath"`
	// AccessTokenDur   int    `yaml:"accessTokenDur"`  // in min
	// RefreeshTokenDur int    `yaml:"refreshTokenDur"` // in min
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
}

func prepare(configPath string) error {
	sc := ServerConfig{
		Port:       8000,
		UploadPath: "/home/kang/Downloads/uploads",
	}

	pg := PostgresConfig{Host: "localhost",
		Port:         5432,
		User:         "kang",
		Password:     "",
		DatabaseName: "kmem",
		SslMode:      "disable",
	}

	conf := Config{
		Server:   sc,
		Postgres: pg,
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

	pgPass := os.Getenv("POSTGRES_PASSWORD")
	if len(pgPass) == 0 {
		return nil, fmt.Errorf("failed to get postgres password")
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecret) == 0 {
		return nil, fmt.Errorf("failed to get jwt secret key")
	}

	conf.Postgres.Password = pgPass
	conf.Server.JwtSecret = jwtSecret

	return &conf, nil
}

func (c *Config) ServerPort() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

func (c *Config) PostgresConnStr() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DatabaseName,
		c.Postgres.SslMode,
	)
}

func (c *Config) JwtSecretKey() string {
	return c.Server.JwtSecret
}

func (c *Config) UploadPath() string {
	return c.Server.UploadPath
}
