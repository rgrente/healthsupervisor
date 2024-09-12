package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RemoteSupervisors []map[string]interface{} "yaml:\"remoteSupervisors\""
	Probes            []map[string]interface{} "yaml:\"probes\""
	Rules             []map[string]interface{} "yaml:\"rules\""
	Hooks             []map[string]interface{} "yaml:\"hooks\""
}

// NewConfig reads a YAML config file and unmarshals it into a Config struct
func NewConfig(configFilePath string) (*Config, error) {
	// Open the file
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Initialize the Config struct
	var c Config

	// Unmarshal YAML content into the Config struct
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	// Return the Config struct

	return &c, nil
}
