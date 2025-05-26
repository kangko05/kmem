package tests

import (
	"kmem/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configPath := "../config.yml"

	t.Run("test load", func(t *testing.T) {
		conf, err := config.Load(configPath)
		assert.Nil(t, err)
		assert.Equal(t, ":8000", conf.ServerPort())
	})
}
