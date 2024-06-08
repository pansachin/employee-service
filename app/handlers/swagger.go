package handlers

import (
	"github.com/pansachin/employee-service/pkg/api"
)

// swagger:response errorResponse400
type _ struct {
	// in:body
	Body struct {
		// Bad Request
		//
		// example: false
		Success bool `json:"success"`
		// Timestamp
		//
		// example: 1639237536
		Timestamp int64             `json:"timestamp"`
		Errors    api.ErrorResponse `json:"errors"`
	}
}

// swagger:response errorResponse404
type _ struct {
	// in:body
	Body struct {
		// Not Found
		//
		// example: false
		Success bool `json:"success"`
		// Timestamp
		//
		// example: 1639237536
		Timestamp int64 `json:"timestamp"`
		// example: {"error": "employee not found"}
		Errors map[string]string `json:"errors"`
	}
}
