// Package dbtest for testing db
package dbtest

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/pansachin/employee-service/pkg/database"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

var (
	UnitDbConfig = database.Config{
		Type:       "mysql",
		User:       "root",
		Password:   "root",
		Host:       "localhost",
		Port:       7801,
		Name:       "employee",
		DisableTLS: true,
	}
)

// NewUnit for initializing unit tests
func NewUnit(t *testing.T) (*slog.Logger, *sqlx.DB, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Open(UnitDbConfig)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	if err := database.StatusCheck(ctx, db); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	t.Log("Creating test database ...")

	// Make sure we have sufficient permission for the db user
	if _, err := db.ExecContext(context.Background(), "CREATE DATABASE IF NOT EXISTS test_db"); err != nil {
		t.Fatalf("dropping database test_db: %v", err)
	}

	t.Log("Test database ready")

	_ = db.Close()

	testDbConfig := UnitDbConfig
	testDbConfig.Name = "test_db"
	db, err = database.Open(testDbConfig)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Ready for testing ...")

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		_ = db.Close()

		_ = writer.Flush()
		log.Info("******************** LOGS ********************")
		log.Info(buf.String())
		log.Info("******************** LOGS ********************")
	}

	return log, db, teardown
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// BoolPointer is a helper to get a *bool from a bool.
// because we normally don't want to deal with pointers to basic types, but it's
// useful in some tests.
func BoolPointer(b bool) *bool {
	return &b
}
