// Package logger for logger functions
package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config holds the is the logger instance configuratoin.
type Config struct {
	// Writer is the writer to write the logs to.
	// It will set the writer to the provided value.
	//
	// If not provided, it will default to os.Stdout.
	Writer io.Writer

	// Source is the source information flag.
	// It will set AddSource to true.
	//
	// If not provided, it will default to true.
	Source bool

	// Debug sets the level of logging
	//
	// By default it's false
	Debug bool

	// json decideds on in which format log will be printed
	//
	// Default it's sets to false and print's logs in text format.
	Json bool
}

// NewLogger return a new instance of logger initalized based on passed config.
func NewLogger(config *Config, appname string, appversion string) (*slog.Logger, error) {
	// Fetch app version if not found
	// When running locally, the VERSION file is in the root of the project
	if appversion == "" {
		// Ignore errors, this is optional
		version, _ := os.ReadFile("../../VERSION")
		if version != nil {
			appversion = strings.TrimSpace(string(version))
		}
	}

	if config.Json {
		return NewJSONLogger(config, appname, appversion)
	}
	return NewTextLogger(config, appname, appversion)
}

// NewTextLogger returns a text logger
// For local use
func NewTextLogger(config *Config, appname string, appversion string) (*slog.Logger, error) {
	// Default is INFO
	var loglevel slog.Level
	if config.Debug {
		loglevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(
		config.Writer,
		&slog.HandlerOptions{
			Level:     loglevel,
			AddSource: config.Source,
		}).WithAttrs([]slog.Attr{
		slog.String("service", appname),
		slog.String("version", appversion),
	})

	return slog.New(handler), nil
}

// NewJSONLogger retrns a JSON logger
// Use in development and production environments - GCP environment
func NewJSONLogger(config *Config, appname string, appversion string) (*slog.Logger, error) {
	// Default is INFO
	var loglevel slog.Level
	if config.Debug {
		loglevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(
		config.Writer,
		&slog.HandlerOptions{
			Level:     loglevel,
			AddSource: config.Source,
		}).WithAttrs([]slog.Attr{
		slog.String("service", appname),
		slog.String("version", appversion),
	})

	return slog.New(handler), nil
}
