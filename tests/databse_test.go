package tests

import (
	"kmem/internal/config"
	"kmem/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	conf, err := config.Load("../config.yml")
	assert.Nil(t, err)

	t.Run("testConnection", func(t *testing.T) {
		pg, err := database.Connect(t.Context(), &conf.Postgres)
		if err != nil {
			t.Fatal(err)
		}
		defer pg.Close()
	})
}

// func TestQuery(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	assert.Nil(t, err)
//
// 	t.Run("test query user", func(t *testing.T) {
// 		ctx, cancel := context.WithCancel(context.Background())
//
// 		pg, err := database.Connect(ctx)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer pg.Close()
//
// 		user, err := query.QueryUser(pg, "guest")
//
// 		assert.Nil(t, err)
// 		assert.Equal(t, "guest", user.Username)
// 		assert.True(t, utils.CheckPasswordHash("guestpass", user.Password))
//
// 		cancel()
// 	})
// }
