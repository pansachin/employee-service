package config

// Config is the configuration for the parser.
type Config struct {
	paths         []string
	serviceConfig ServiceConfig
}

// NewConfig creates a new Config.
func NewConfig() *Config {
	return &Config{}
}

// SetConfigPaths sets the paths for the configuration.
// For e.g., the paths can be:
//
//	The first path should be the cloud run configuration in GCP as a mounted secrets.
//	The second path should be the local configuration.
func (c *Config) SetConfigPaths(paths []string) error {
	if len(paths) == 0 {
		return ErrEmptyConfigPath
	}
	c.paths = paths

	return nil
}

func (c *Config) GetServiceConfig() ServiceConfig {
	return c.serviceConfig
}
