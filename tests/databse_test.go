package tests

import (
	"context"
	"kmem/internal/database"
	"kmem/internal/database/query"
	"kmem/internal/utils"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	err := godotenv.Load("../.env")
	assert.Nil(t, err)

	t.Run("testConnection", func(t *testing.T) {
		pg, err := database.Connect(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		defer pg.Close()
	})
}

func TestQuery(t *testing.T) {
	err := godotenv.Load("../.env")
	assert.Nil(t, err)

	t.Run("test query user", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		pg, err := database.Connect(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer pg.Close()

		user, err := query.QueryUser(pg, "guest")

		assert.Nil(t, err)
		assert.Equal(t, "guest", user.Username)
		assert.True(t, utils.CheckPasswordHash("guestpass", user.Password))

		cancel()
	})
}
