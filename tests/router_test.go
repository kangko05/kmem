package tests

import (
	"bytes"
	"encoding/json"
	"kmem/internal/models"
	"kmem/internal/router"
	"kmem/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	router := router.Setup(testDB, testConfig, testQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong\n", w.Body.String())
}

func TestSignupSuccess(t *testing.T) {
	cleanupTables(t)

	router := router.Setup(testDB, testConfig, testQueue)
	testUser := models.User{
		Username: "testuser",
		Password: "testpassword123",
	}

	wb, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewReader(wb))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	u, err := testDB.QueryUser(testUser.Username)
	assert.Nil(t, err)
	assert.Equal(t, testUser.Username, u.Username)
	assert.True(t, utils.CheckPasswordHash(u.Password, testUser.Password))
}

func TestSignupValidation(t *testing.T) {
	cleanupTables(t)

	router := router.Setup(testDB, testConfig, testQueue)

	tests := []struct {
		name string
		user models.User
		want int
	}{
		{"short username", models.User{Username: "ab", Password: "testpassword123"}, http.StatusBadRequest},
		{"short password", models.User{Username: "testuser", Password: "123"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb, _ := json.Marshal(tt.user)
			req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewReader(wb))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	cleanupTables(t)

	router := router.Setup(testDB, testConfig, testQueue)
	testUser := models.User{
		Username: "testuser",
		Password: "testpassword123",
	}

	err := testDB.InsertUser(testUser)
	assert.Nil(t, err)

	wb, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(wb))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	cookieNames := make(map[string]bool)
	for _, cookie := range cookies {
		cookieNames[cookie.Name] = true
	}
	assert.True(t, cookieNames[utils.ACCESS_TOKEN_KEY])
	assert.True(t, cookieNames[utils.REFRESH_TOKEN_KEY])
}

func TestLoginWrongCredentials(t *testing.T) {
	cleanupTables(t)

	router := router.Setup(testDB, testConfig, testQueue)

	wrongUser := models.User{
		Username: "wronguser",
		Password: "wrongpassword",
	}

	wb, _ := json.Marshal(wrongUser)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(wb))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
