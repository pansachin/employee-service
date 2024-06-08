package config

import "errors"

var (
	ErrEmptyConfigData   = errors.New("empty config data")
	ErrEmptyConfigPath   = errors.New("empty config path")
	ErrNoConfigFileFound = errors.New("no config file found in the given paths")
)
