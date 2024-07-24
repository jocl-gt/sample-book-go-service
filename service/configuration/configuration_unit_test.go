package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test-config-*.yaml")
	assert.NoError(t, err)
	defer tmpfile.Close()

	yamlConfig := `
db_username: test_user
db_password: test_password
db_hostname: localhost
db_port: 5432
db_database_name: test_db
db_ssl_mode: disable
server_port: 8080
cors_allowed_origins:
  - http://localhost
`

	_, err = tmpfile.WriteString(yamlConfig)
	assert.NoError(t, err)

	config := LoadConfiguration(tmpfile.Name())

	expectedConfig := Configuration{
		DBUsername:     "test_user",
		DBPassword:     "test_password",
		DBHostname:     "localhost",
		DBPort:         5432,
		DBDatabaseName: "test_db",
		DBSSLMode:      "disable",
		ServerPort:     8080,
		CORSOrigins: []string{
			"http://localhost",
		},
	}
	assert.Equal(t, expectedConfig, config)
}
