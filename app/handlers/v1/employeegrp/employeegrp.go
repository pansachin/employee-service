// Package Employeegrp for requestion source handler functions
package employeegrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pansachin/employee-service/models/employee"
	"github.com/pansachin/employee-service/pkg/api"
	"github.com/pansachin/employee-service/pkg/database"
)

// Handlers manages the set of employee endpoints.
type Handlers struct {
	Employee employee.Core
}

// Create a new employee record
//
// # Create a new Employee record
//
// ---
// - application/json
// responses:
//
//	  "200":
//		   "$ref": "#/responses/EmployeeRes"
//	  "400":
//		   "$ref": "#/responses/errorResponse400"
//	  "404":
//		   "$ref": "#/responses/errorResponse404"
//
//swagger:operation POST /employee Employee EmployeeCreate
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	nes := employee.NewEmployee{}
	if err := api.Decode(r, &nes); err != nil {
		return api.NewRequestError(err, http.StatusBadRequest)
	}

	now := time.Now().UTC()

	data, err := h.Employee.Create(ctx, nes, now)
	if err != nil {
		switch {
		case errors.Is(err, employee.ErrInvalidID):
			return api.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, employee.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("something went wrong, ERRPR: [%w]", err)
		}
	}

	return api.Respond(ctx, w, []employee.Employee{data}, http.StatusOK)
}

// Query all the Employee records
//
// swagger:operation GET /employee Employee EmployeeQuery
//
// # This is the summary for listing Employeees
//
// ---
// produces:
// - application/json
// responses:
//
//	  "200":
//		   "$ref": "#/responses/EmployeeRes"
//	  "404":
//		   "$ref": "#/responses/errorResponse404"
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	pagi, err := database.PaginationParams(r)
	if err != nil {
		return err
	}

	rs, err := h.Employee.Query(ctx, pagi)
	if err != nil {
		switch {
		case errors.Is(err, employee.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("unable to query for Employee: %w", err)
		}
	}

	return api.Respond(ctx, w, rs, http.StatusOK)
}

// QueryByID from an individual id
//
// swagger:operation GET /employee/{id} Employee EmployeeQueryById
//
// # Getting a single Employee by ID
//
// ---
// produces:
// - application/json
// responses:
//
//	  "200":
//		   "$ref": "#/responses/EmployeeRes"
//	  "404":
//		   "$ref": "#/responses/errorResponse404"
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := api.Param(r, "id")

	rs, err := h.Employee.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, employee.ErrInvalidID):
			return api.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, employee.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("employee id[%s]: %w", id, err)
		}
	}

	return api.Respond(ctx, w, []employee.Employee{rs}, http.StatusOK)
}

// Delete from an individual id
//
// swagger:operation DELETE /employee/{id} Employee EmployeeDelete
//
// # Delete a single Employee by ID
//
// ---
// produces:
// - application/json
// responses:
//
//	  "200":
//		   "$ref": "#/responses/EmployeeRes"
//	  "404":
//		   "$ref": "#/responses/errorResponse404"
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := api.Param(r, "id")

	now := time.Now().UTC()

	err := h.Employee.Delete(ctx, id, now)
	if err != nil {
		switch {
		case errors.Is(err, employee.ErrInvalidID):
			return api.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, employee.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("employee id[%s]: %w", id, err)
		}
	}

	return api.Respond(ctx, w, nil, http.StatusOK)
}

// Update from an individual id
//
// swagger:operation PATCH /employee/{id} Employee EmployeeUpdate
//
// # Update a single Employee by ID
//
// ---
// produces:
// - application/json
// responses:
//
//	  "200":
//		   "$ref": "#/responses/EmployeeRes"
//	  "404":
//		   "$ref": "#/responses/errorResponse404"
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := api.Param(r, "id")

	ues := employee.UpdateEmployee{}
	if err := api.Decode(r, &ues); err != nil {
		return api.NewRequestError(err, http.StatusBadRequest)
	}

	now := time.Now().UTC()

	err := h.Employee.Update(ctx, id, ues, now)
	if err != nil {
		switch {
		case errors.Is(err, employee.ErrInvalidID):
			return api.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, employee.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("employee id[%s]: %w", id, err)
		}
	}

	return api.Respond(ctx, w, nil, http.StatusOK)
}

// UnDelete from an individual id
//
// swagger:operation PATCH /employee/undelete/{id} Employee EmployeeUnDelete
//
// # UnDelete a single Employee by ID
//
// ---
// produces:
// - application/json
// responses:
//
//	  "200":
//		   "$ref": "#/responses/EmployeeRes"
//	  "404":
//		   "$ref": "#/responses/errorResponse404"
func (h Handlers) UnDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := api.Param(r, "id")

	now := time.Now().UTC()

	err := h.Employee.UnDelete(ctx, id, now)
	if err != nil {
		switch {
		case errors.Is(err, employee.ErrInvalidID):
			return api.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, employee.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("employee id[%s]: %w", id, err)
		}
	}

	return api.Respond(ctx, w, nil, http.StatusOK)
}
