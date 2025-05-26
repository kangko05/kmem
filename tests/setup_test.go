package tests

import (
	"kmem/internal/config"
	"kmem/internal/db"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testDB *db.Postgres
var testConfig *config.Config

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Failed to load .env file")
	}

	conf, err := config.Load("../config.yml")
	if err != nil {
		panic("Failed to load config")
	}
	testConfig = conf

	pg, err := db.Connect(conf)
	if err != nil {
		panic("Failed to connect to test database")
	}
	testDB = pg

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func cleanupTables(t *testing.T) {
	err := testDB.Exec("TRUNCATE TABLE users CASCADE")
	assert.Nil(t, err)
}
