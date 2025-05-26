package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	assert.Nil(t, testDB.Ping())
}
