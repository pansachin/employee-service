// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/pansachin/employee-service/app/handlers/v1/employeegrp"
	"github.com/pansachin/employee-service/models/employee"
	"github.com/pansachin/employee-service/pkg/api"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *slog.Logger
	DB    *sqlx.DB
	RWMux *sync.RWMutex
}

// Routes binds all the version 1 routes.
func Routes(router *api.API, cfg Config) {
	// -------------------------------------------------------------------
	// Requesting Sources
	// -------------------------------------------------------------------
	rs := employeegrp.Handlers{
		Employee: employee.NewCore(cfg.Log, cfg.DB, cfg.RWMux),
	}
	router.Handle(http.MethodPost, "/v1/employee", rs.Create)
	router.Handle(http.MethodGet, "/v1/employee", rs.Query)
	router.Handle(http.MethodGet, "/v1/employee/{id}", rs.QueryByID)
	router.Handle(http.MethodPatch, "/v1/employee/{id}", rs.Update)
	router.Handle(http.MethodDelete, "/v1/employee/{id}", rs.Delete)
	router.Handle(http.MethodPatch, "/v1/employee/undelete/{id}", rs.UnDelete)

	// -------------------------------------------------------------------
	// Add in the Teapot
	// -------------------------------------------------------------------
	router.Handle(http.MethodGet, "/v1/teapot", Teapot)
}

// Teapot ... because everyone needs to return a 418 at some point
func Teapot(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	lyrics := `I'm a little teapot, Short and stout,
Here is my handle. Here is my spout.
When I get all steamed up, Hear me shout,
Tip me over and pour me out!`

	type o struct {
		Lyrics string `json:"lyrics"`
	}
	output := []o{{Lyrics: lyrics}}

	return api.Respond(ctx, w, output, http.StatusTeapot)
}
