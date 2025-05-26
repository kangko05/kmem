package tests

import (
	"kmem/internal/config"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Nil(t, godotenv.Load("../.env"))

	configPath := "../config.yml"

	t.Run("test load", func(t *testing.T) {
		conf, err := config.Load(configPath)
		assert.Nil(t, err)
		assert.Equal(t, ":8000", conf.ServerPort())
	})
}
