package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

/*
	The configuration package provides the functionality to read a YAML file containing
	the information required to load the application.
*/

type Configuration struct {
	DBUsername     string   `yaml:"db_username"`
	DBPassword     string   `yaml:"db_password"`
	DBHostname     string   `yaml:"db_hostname"`
	DBPort         int      `yaml:"db_port"`
	DBDatabaseName string   `yaml:"db_database_name"`
	DBSSLMode      string   `yaml:"db_ssl_mode"`
	ServerPort     int      `yaml:"server_port"`
	CORSOrigins    []string `yaml:"cors_allowed_origins"`
}

func LoadConfiguration(configurationFilePath string) Configuration {
	var config Configuration

	if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
		configurationFilePath = envConfigFile
	}

	configFile, err := os.ReadFile(configurationFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
