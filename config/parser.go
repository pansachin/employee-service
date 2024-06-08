package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Parse reads the configuration file from the given paths.
func (c *Config) Parse(filename string) error {
	// check if the paths are empty
	if len(c.paths) == 0 {
		return ErrEmptyConfigPath
	}

	// check if the file exists
	for _, path := range c.paths {
		srvConfig := ServiceConfig{}
		configPath := filepath.Join(path, filename)

		// Check if the file exists.
		_, err := os.Stat(configPath)
		if err != nil {
			continue
		}

		// Read the file
		cfgData, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}

		// Check if the file is empty
		if len(cfgData) == 0 {
			return ErrEmptyConfigData
		}

		// Unmarshal the file
		err = yaml.Unmarshal(cfgData, &srvConfig)
		if err != nil {
			return err
		}
		c.serviceConfig = srvConfig
		return nil
	}

	return ErrNoConfigFileFound
}
