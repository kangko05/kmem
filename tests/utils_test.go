package tests

import (
	"kmem/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {

	const JWT_SECRET = "TODO: SECRET KEY MUST BE HANDLED WITH CARE"

	t.Run("create jwt", func(t *testing.T) {
		tokenStr, err := utils.CreateJwt(time.Hour, "testuser", JWT_SECRET)
		assert.Nil(t, err)
		assert.NotEqual(t, len(tokenStr), 0)
	})

	t.Run("parse valid jwt", func(t *testing.T) {
		tokenStr, err := utils.CreateJwt(time.Hour, "testuser", JWT_SECRET)
		assert.Nil(t, err)
		assert.NotEqual(t, len(tokenStr), 0)

		_, _, err = utils.PasrseJwt(JWT_SECRET, tokenStr)
		assert.Nil(t, err)
	})

	t.Run("extract username from token", func(t *testing.T) {
		tokenStr, err := utils.CreateJwt(time.Hour, "testuser", JWT_SECRET)
		assert.Nil(t, err)
		assert.NotEqual(t, len(tokenStr), 0)

		_, claim, err := utils.PasrseJwt(tokenStr, JWT_SECRET)
		assert.Nil(t, err)

		v, ok := claim[utils.USERNAME_KEY]
		assert.True(t, ok)
		assert.Equal(t, "testuser", v)
	})
}
