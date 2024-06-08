// Package employee Golang sample Service Microservice API.
//
// Terms Of Service:
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//	Schemes: http, https
//	Host: localhost
//	Description: This is a base sample service
//	BasePath: /v1
//	Version: 0.0.1
//	License: MIT http://opensource.org/licenses/MIT
//	Contact: Sachin Prasad <prasadsachin214@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	//nolint:all

	"github.com/pansachin/employee-service/app/handlers"
	"github.com/pansachin/employee-service/config"
	"github.com/pansachin/employee-service/pkg/database"
	"github.com/pansachin/employee-service/pkg/logger"
)

const (
	appName = "employee-service"
)

var (
	appVersionLDFlag        string
	appBuildTimestampLDFlag string
)

func main() {

	// -------------------------------------------------------------------
	// Logger
	// -------------------------------------------------------------------
	log, err := logger.NewLogger(&logger.Config{
		Writer: os.Stdout,
		Json:   true,
	},
		appName,
		appVersionLDFlag,
	)

	if err != nil {
		fmt.Println("error initializing production logger")
		os.Exit(1)
	}

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Error("startup failure", slog.Any("ERROR", err))

		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	var err error
	configYMLFile := "employee-service.config.yml"

	c := config.NewConfig()
	// Read the config files as per the priority.
	err = c.SetConfigPaths([]string{
		// Reads cloud function configuration in GCP as a mounted secrets.
		// Same path should be used while provisioning the secret to CF.
		"/service/config",
		// Reads local system configurations.
		".",
	})
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err = c.Parse(configYMLFile); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	srvCfg := c.GetServiceConfig()

	// -------------------------------------------------------------------
	// Startup Details
	// -------------------------------------------------------------------
	log.Info("startup", "binary build time", appBuildTimestampLDFlag)

	// Reinitalize logger with proper configuration
	log, err = logger.NewLogger(&logger.Config{
		Writer: os.Stdout,
		Source: srvCfg.Log.Source,
		Debug:  srvCfg.Log.Debug,
		Json:   srvCfg.Log.JSON,
	},
		appName,
		appVersionLDFlag,
	)
	if err != nil {
		log.Error("configuration failure", "section", ".startup.configs", slog.Any("ERROR", err))
	}

	fmt.Printf("employee-service-configs: %%#v: %#v\n", srvCfg)

	// -------------------------------------------------------------------
	// Databases
	// -------------------------------------------------------------------
	log.Info("startup.db", "status", "initializing DBs")

	db, err := database.Open(database.Config{
		Type:         srvCfg.Db.Type,
		User:         srvCfg.Db.User,
		Password:     srvCfg.Db.Password,
		Host:         srvCfg.Db.Host,
		Port:         srvCfg.Db.Port,
		Name:         srvCfg.Db.DbName,
		MaxIdleConns: srvCfg.Db.MaxIdleConns,
		MaxOpenConns: srvCfg.Db.MaxOpenConns,
		DisableTLS:   srvCfg.Db.DisableTLS,
	})
	if err != nil {
		panic(fmt.Errorf("connecting to db: %w", err))
	}
	defer func() {
		log.Info("shutdown", "status", "stopping db", "host", srvCfg.Db.Host)
		_ = db.Close()
	}()

	// -------------------------------------------------------------------
	// RWMux for lock DBs in transaction mode (deadlocks = yuck)
	// -------------------------------------------------------------------
	log.Info("startup.remux", "status", "created")
	rwmux := &sync.RWMutex{}

	// -------------------------------------------------------------------
	// Initialize API
	// -------------------------------------------------------------------
	log.Info("startup.api", "status", "initializing API")

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Env:     srvCfg.App.Env,
		Log:     log,
		DB:      db,
		RWMux:   rwmux,
		Headers: srvCfg.App.EnforceHeaders,
	})

	// -------------------------------------------------------------------
	// New Channels
	// -------------------------------------------------------------------

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	subscribeErrors := make(chan error, 1)

	apiHost := fmt.Sprintf("%s:%s", srvCfg.Web.APIHost, srvCfg.Web.APIPort)
	api := http.Server{
		Addr:              apiHost,
		Handler:           apiMux,
		ReadHeaderTimeout: srvCfg.Web.ReadHeaderTimeout,
		ReadTimeout:       srvCfg.Web.ReadTimeout,
		WriteTimeout:      srvCfg.Web.WriteTimeout,
		IdleTimeout:       srvCfg.Web.IdleTimeout,
		MaxHeaderBytes:    srvCfg.Web.MaxHeaderBytes,
	}
	// TODO: Push this in with a new Interface for a logger w/ io.Writer
	// https://stackoverflow.com/questions/52294334/net-http-set-custom-logger
	//api.ErrorLog = stdLibLog.New(log, "", 0)
	//api.TLSConfig.BuildNameToCertificate()

	// -------------------------------------------------------------------
	// Starting the API
	// -------------------------------------------------------------------

	// Start the service listening for api requests.
	go func() {
		log.Info("startup.api", "status", "api router started", "host", api.Addr)

		log.Debug("API STARTED", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------
	// Shutdown
	// -------------------------------------------------------------------

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case err := <-subscribeErrors:
		return fmt.Errorf("subscriber error: %w", err)

	case sig := <-shutdown:
		log.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info("shutdown", "status", "shutdown completed", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), srvCfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			_ = api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
