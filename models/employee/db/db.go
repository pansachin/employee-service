// Package db for database functions
package db

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/pansachin/employee-service/pkg/database"
)

// Store holds details for basic database needs
type Store struct {
	log          *slog.Logger
	tr           database.Transactor
	db           sqlx.ExtContext
	rwmux        *sync.RWMutex
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *slog.Logger, db *sqlx.DB, rwmux *sync.RWMutex) Store {
	return Store{
		log:   log,
		tr:    db,
		db:    db,
		rwmux: rwmux,
	}
}

// WithinTran runs passes function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		return fn(s.db)
	}
	s.rwmux.Lock()
	err := database.WithinTran(ctx, s.log, s.tr, fn)
	s.rwmux.Unlock()

	return err
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// -----------------------------------------------------------------------
// Database Query Repository
// -----------------------------------------------------------------------

// Create inserts a new requesting into the database.
func (s Store) Create(ctx context.Context, rs Employee) (database.DBResults, error) {
	const q = `
	INSERT INTO employee
		(name, position, created_on, updated_on)
	VALUES
		(:name, :position, :created_on, :updated_on)`

	res, err := database.NamedExecContext(ctx, s.log, s.db, q, rs)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return database.DBResults{}, database.NewError(database.ErrDBDuplicatedEntry, http.StatusConflict)
		}
		return database.DBResults{}, fmt.Errorf("inserting employee: %w", err)
	}

	return res, nil
}

// Update replaces a employee record in the database.
func (s Store) Update(ctx context.Context, rs Employee) (database.DBResults, error) {
	const q = `
	UPDATE
		employee
	SET 
		name = :name,
		position = :position,
		updated_on = :updated_on
	WHERE
		id = :id`

	res, err := database.NamedExecContext(ctx, s.log, s.db, q, rs)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return database.DBResults{}, database.NewError(database.ErrDBDuplicatedEntry, http.StatusConflict)
		}
		return database.DBResults{}, fmt.Errorf("updating Employee ID[%s]: %w", rs.ID, err)
	}

	return res, nil
}

// Delete removes a employee from the database.
func (s Store) Delete(ctx context.Context, id string, now time.Time) (database.DBResults, error) {
	data := struct {
		ID        string    `db:"id"`
		DeletedOn time.Time `db:"deleted_on"`
	}{
		ID:        id,
		DeletedOn: now,
	}

	const q = `
	UPDATE
		employee
	SET
		deleted_on = :deleted_on
	WHERE
		id = :id`

	res, err := database.NamedExecContext(ctx, s.log, s.db, q, data)
	if err != nil {
		return database.DBResults{}, fmt.Errorf("deleting employee id[%s]: %w", id, err)
	}

	return res, nil
}

// Query retrieves a list of existing employee from the database.
func (s Store) Query(ctx context.Context, pagi database.Pagination) ([]Employee, error) {
	q := database.PaginationQuery(pagi, `
	SELECT
		id,
	    name,
	    position,
	    created_on,
	    updated_on,
	    deleted_on
	FROM
		employee
	WHERE
		deleted_on is null
	ORDER BY
		:sort :direction,
		id :direction
	LIMIT
		:page,:per_page`)

	// Slice to hold results
	var res []Employee
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, pagi, &res); err != nil {
		return nil, fmt.Errorf("selecting employee: %w", err)
	}

	return res, nil
}

// QueryByID retrieves a list of existing requesting sources from the database.
func (s Store) QueryByID(ctx context.Context, id string) (Employee, error) {
	data := struct {
		ID string `db:"id"`
	}{ID: id}

	const q = `
	SELECT
		id,
		name,
		position,
		created_on,
		updated_on,
		deleted_on
	FROM
		employee
	WHERE
		id = :id
		and deleted_on is null`

	// Slice to hold results
	var res Employee
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &res); err != nil {
		return Employee{}, fmt.Errorf("selecting by id[%q]: %w", id, err)
	}

	return res, nil
}

// UnDelete restores a deleted employee from the database.
func (s Store) UnDelete(ctx context.Context, id string, now time.Time) (database.DBResults, error) {
	data := struct {
		ID        string    `db:"id"`
		UpdatedOn time.Time `db:"updated_on"`
		DeletedOn time.Time `db:"deleted_on"`
	}{
		ID:        id,
		UpdatedOn: now,
		DeletedOn: time.Time{},
	}

	const q = `
	UPDATE
		employee
	SET
		updated_on = :updated_on,
		deleted_on = null
	WHERE
		id = :id`

	res, err := database.NamedExecContext(ctx, s.log, s.db, q, data)
	if err != nil {
		return database.DBResults{}, fmt.Errorf("restoring employee id[%s]: %w", id, err)
	}

	return res, nil
}
