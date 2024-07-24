package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerInstance(t *testing.T) {
	configurationFilePath := "../config.integration.yaml"
	server, err := NewServerInstance(configurationFilePath)
	assert.NoError(t, err)
	server.Start()
	server.Close()
}
