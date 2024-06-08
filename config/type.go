package config

import "time"

// ServiceConfig is the configuration for the service.
type ServiceConfig struct {
	App App `yaml:"app"`
	Web Web `yaml:"web"`
	Log Log `yaml:"log"`
	Db  Db  `yaml:"db"`
}

// App is the configuration for the app.
type App struct {
	Name           string `yaml:"name"`
	Env            string `yaml:"env"`
	EnforceHeaders bool   `yaml:"enforceHeaders"`
	TLS            bool   `yaml:"tls"`
	Function       string `yaml:"function"`
}

// Web is the configuration for the web.
type Web struct {
	MaxHeaderBytes    int           `yaml:"maxHeaderBytes"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
	ReadTimeout       time.Duration `yaml:"readTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdownTimeout"`
	APIHost           string        `yaml:"apiHost"`
	APIPort           string        `yaml:"apiPort"`
}

// Log is the configuration for the log.
type Log struct {
	Debug  bool `yaml:"debug"`
	JSON   bool `yaml:"json"`
	Source bool `yaml:"source"`
}

// Db is the configuration for the db.
type Db struct {
	Type         string `yaml:"type"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DbName       string `yaml:"dbName"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	DisableTLS   bool   `yaml:"disableTLS"`
}
