package employeegrp

import "github.com/pansachin/employee-service/models/employee"

// swagger:response EmployeeRes
type _ struct {
	// in:body
	Body struct {
		// Success
		//
		Success bool `json:"success"`
		// Timestamp
		//
		// example: 1639237536
		Timestamp int64 `json:"timestamp"`
		// Data
		// in: body
		Data []employee.Employee `json:"data"`
	}
}

// swagger:parameters EmployeeQueryById EmployeeDelete EmployeeUpdate EmployeeUnDelete
type _ struct {
	// Employee ID
	//
	// in: path
	// required: true
	// enum: 1
	// type: integer
	ID string `json:"id"`
}

// swagger:parameters EmployeeCreate
type _ struct {
	// The body to create a employee
	// in:body
	// required: true
	Body employee.NewEmployee
}

// swagger:parameters EmployeeUpdate
type _ struct {
	// The body to update a employee
	// in:body
	// required: true
	Body employee.UpdateEmployee
}
