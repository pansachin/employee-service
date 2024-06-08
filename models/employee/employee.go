// Package Employee for employee handler functions
package employee

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/pansachin/employee-service/models/employee/db"
	"github.com/pansachin/employee-service/pkg/database"
	"github.com/pansachin/employee-service/pkg/validate"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("employee not found")
	ErrInvalidID    = errors.New("ID is not in its proper form")
	ErrInvalidAlias = errors.New("alias is not in its proper form")
)

// Core manages the set of APIs for employee access
type Core struct {
	store db.Store
}

// NewCore constructs a core for employee api access.
func NewCore(log *slog.Logger, sqlxDB *sqlx.DB, rwmux *sync.RWMutex) Core {
	return Core{
		store: db.NewStore(log, sqlxDB, rwmux),
	}
}

// -----------------------------------------------------------------------
// CRUD Methods
// -----------------------------------------------------------------------

// Create inserts a new employee into the database
func (c Core) Create(ctx context.Context, rs NewEmployee, now time.Time) (Employee, error) {
	if err := validate.Check(rs); err != nil {
		return Employee{}, fmt.Errorf("validating data: %w", err)
	}

	dbRS := db.Employee{
		Name:      strings.TrimSpace(rs.Name),
		Position:  strings.TrimSpace(rs.Position),
		CreatedOn: now,
		UpdatedOn: now,
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		res, err := c.store.Tran(tx).Create(ctx, dbRS)
		if err != nil {
			return err
		}
		dbRS.ID = fmt.Sprintf("%d", res.LastInsertID)
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return Employee{}, fmt.Errorf("tran: %w", err)
	}

	return toEmployee(dbRS), nil
}

// Update replaces a employee document in the database.
func (c Core) Update(ctx context.Context, id string, urs UpdateEmployee, now time.Time) error {
	if err := validate.Check(urs); err != nil {
		return err
	}
	if err := validate.CheckID(id); err != nil {
		return ErrInvalidID
	}

	dbRS, err := c.store.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating employee id[%s]: %w", id, err)
	}

	isEmpty := true
	if urs.Position != nil {
		dbRS.Position = strings.TrimSpace(*urs.Position)
		isEmpty = false
	}
	// No changes were made - don't touch the DB
	if isEmpty {
		return nil
	}
	dbRS.UpdatedOn = now

	_, err = c.store.Update(ctx, dbRS)
	if err != nil {
		return fmt.Errorf("update id[%s]: %w", id, err)
	}

	return nil
}

// Delete removes a employee from the database.
func (c Core) Delete(ctx context.Context, id string, now time.Time) error {
	if err := validate.CheckID(id); err != nil {
		return ErrInvalidID
	}

	_, err := c.store.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("undeleting employee id[%s]: %w", id, err)
	}

	_, err = c.store.Delete(ctx, id, now)
	if err != nil {
		return fmt.Errorf("delete id[%s]: %w", id, err)
	}

	return nil
}

// Query retrieves a list of existing records from the database
func (c Core) Query(ctx context.Context, pagi database.Pagination) ([]Employee, error) {
	res, err := c.store.Query(ctx, pagi)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toEmployeeSlice(res), nil
}

// QueryByID retrieves a single records from the database by id
func (c Core) QueryByID(ctx context.Context, id string) (Employee, error) {
	if err := validate.CheckID(id); err != nil {
		return Employee{}, ErrInvalidID
	}

	res, err := c.store.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Employee{}, ErrNotFound
		}
		return Employee{}, fmt.Errorf("query: %w", err)
	}

	return toEmployee(res), nil
}

// UnDelete restore a deleted employee from the database.
func (c Core) UnDelete(ctx context.Context, id string, now time.Time) error {
	if err := validate.CheckID(id); err != nil {
		return ErrInvalidID
	}

	_, err := c.store.UnDelete(ctx, id, now)
	if err != nil {
		return fmt.Errorf("employee id[%s]: %w", id, err)
	}

	return nil
}
