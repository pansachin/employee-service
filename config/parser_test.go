package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigParse(t *testing.T) {
	var (
		c   *Config
		err error
	)
	c = NewConfig()

	// Valid config path
	err = c.SetConfigPaths([]string{"/service/config", "../."})
	require.NoError(t, err)
	err = c.Parse("employee-service.local.yml")
	require.NoError(t, err)
	fmt.Println("Valid config path: passes")

	// Invalid config path
	err = c.SetConfigPaths([]string{"/random/config/path", "."})
	require.NoError(t, err)
	err = c.Parse("employee-service.local.yml")
	assert.Error(t, err, ErrNoConfigFileFound)
	fmt.Println("Invalid config path: passes")

	// Check if the service config is not nil
	err = c.SetConfigPaths([]string{})
	assert.Error(t, err, ErrEmptyConfigPath)
	fmt.Println("Check if the service config is not nil: passes")
}
