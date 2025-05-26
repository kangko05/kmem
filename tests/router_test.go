package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/router"
	"kmem/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// test endpoints
func TestRouter(t *testing.T) {
	assert := assert.New(t)

	err := godotenv.Load("../.env")
	assert.Nil(err)

	conf, err := config.Load("../config.yml")
	assert.Nil(err)

	pg, err := db.Connect(conf)
	assert.Nil(err)

	url := "http://localhost:8000"
	router := router.Setup(pg, conf)

	testUser := models.User{
		Username: "test",
		Password: "testpassword",
	}

	t.Run("test ping", func(t *testing.T) {
		writer := httptest.NewRecorder()

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping", url), nil)
		assert.Nil(err)

		router.ServeHTTP(writer, req)

		assert.Equal(http.StatusOK, writer.Code)
		assert.Equal("pong\n", writer.Body.String())
	})

	// fail test
	t.Run("test signup1", func(t *testing.T) {
		writer := httptest.NewRecorder()

		wb, err := json.Marshal(models.User{
			Username: "a",
			Password: "b",
		})
		assert.Nil(err)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/signup", url), bytes.NewReader(wb))
		assert.Nil(err)

		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(writer, req)

		assert.Equal(http.StatusBadRequest, writer.Code)
	})

	t.Run("test signup2", func(t *testing.T) {
		u, err := pg.QueryUser(testUser.Username)
		assert.Nil(err)

		check := utils.CheckPasswordHash(u.Password, testUser.Password)
		assert.True(check)
	})

	t.Run("test login1", func(t *testing.T) {
		writer := httptest.NewRecorder()

		wb, err := json.Marshal(testUser)
		assert.Nil(err)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/login", url), bytes.NewReader(wb))
		assert.Nil(err)

		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(writer, req)

		for _, cookie := range writer.Result().Cookies() {
			accessOrRefresh := (cookie.Name == utils.ACCESS_TOKEN_KEY) || (cookie.Name == utils.REFRESH_TOKEN_KEY)
			assert.True(accessOrRefresh)
		}

		assert.Equal(http.StatusOK, writer.Code)
	})
}
