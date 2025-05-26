package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kmem/internal/models"
	"kmem/internal/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test endpoints
func TestRouter(t *testing.T) {
	assert := assert.New(t)

	url := "http://localhost:8000"
	router := router.Setup()

	t.Run("test ping", func(t *testing.T) {
		writer := httptest.NewRecorder()

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping", url), nil)
		assert.Nil(err)

		router.ServeHTTP(writer, req)

		assert.Equal(http.StatusOK, writer.Code)
		assert.Equal("pong\n", writer.Body.String())
	})

	t.Run("test signup", func(t *testing.T) {
		writer := httptest.NewRecorder()

		wb, err := json.Marshal(models.User{
			Username: "test",
			Password: "test",
		})
		assert.Nil(err)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/signup", url), bytes.NewReader(wb))
		assert.Nil(err)

		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(writer, req)

		assert.Equal(http.StatusOK, writer.Code)
	})
}
