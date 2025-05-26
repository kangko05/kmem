package tests

import (
	"kmem/internal/config"
	"kmem/internal/db"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	assert.Nil(t, godotenv.Load("../.env"))

	conf, err := config.Load("../config.yml")
	assert.Nil(t, err)

	t.Run("test connection", func(t *testing.T) {
		pg, err := db.Connect(conf)
		assert.Nil(t, err)
		assert.Nil(t, pg.Ping())
		assert.Nil(t, pg.Close())
	})
}
