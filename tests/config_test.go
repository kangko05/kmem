package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	assert.NotNil(t, testConfig)
	assert.Equal(t, ":8000", testConfig.ServerPort())
}
