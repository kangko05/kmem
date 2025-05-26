package tests

import (
	"kmem/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	t.Run("test connection", func(t *testing.T) {
		pg, err := db.Connect()
		assert.Nil(t, err)
		assert.Nil(t, pg.Ping())
		assert.Nil(t, pg.Close())
	})
}
