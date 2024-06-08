package employee

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"github.com/pansachin/employee-service/models/employee/db"
)

// Employee holds the employee information.
//
//swagger:model Employee
type Employee struct {
	// Primary Key
	// example: 1
	ID string `json:"id"`
	// Employee Name
	// example: Sachin Prasad
	Name string `json:"name"`
	// Employee designation
	// example: Senior Software Engineer
	Position string `json:"position"`
	// Database created value
	// example: 2021-05-25T00:53:16.535668Z
	CreatedOn time.Time `json:"created_on"`
	// Database last updated value
	// example: 2021-05-25T00:53:16.535668Z
	UpdatedOn time.Time `json:"updated_on"`
	// Database soft delete value
	// example: 2021-05-25T00:53:16.535668Z
	// swagger:ignore
	DeletedOn *time.Time `json:"deleted_on,omitempty"`
}

// NewEmployee defines the model of adding new employee.
//
//swagger:model NewEmployee
type NewEmployee struct {
	// Name of the employee
	// in: string
	// required: true
	// example: Sachin Prasad
	Name string `json:"name" validate:"required,notblank"`
	// Employee Designamtion
	// in: string
	// example: Senior Software Engineer
	Position string `json:"position"`
}

// UpdateEmployee defines what information may be provided to
// modify an existing Employee. All fields are optional
// so clients can send just the fields they want changed. It uses pointer
// fields so we can differentiate between a field that was not provided
// and a field that was provided as explicitly blank. Normally we do not
// want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
//
//swagger:model UpdateEmployee
type UpdateEmployee struct {
	// Employee Designamtion
	// in: string
	// example: Staff Software Engineer
	Position *string `json:"position"`
}

// =============================================================================

func toEmployee(dbRS db.Employee) Employee {
	p := (*Employee)(unsafe.Pointer(&dbRS))
	return *p
}

func toEmployeeSlice(dbSRs []db.Employee) []Employee {
	rs := make([]Employee, len(dbSRs))
	for i, dbSR := range dbSRs {
		rs[i] = toEmployee(dbSR)
	}
	return rs
}

//------------------------------------------------------------------------
// Fake data generators
//------------------------------------------------------------------------

// GenerateFakeData return an array for NewEmployees
func (nrt NewEmployee) GenerateFakeData(num int) []NewEmployee {
	var data []NewEmployee
	for i := 0; i < num; i++ {
		data = append(data, nrt.fakeData(i+1))
	}
	return data
}

// fakeData creates the fake record
func (nrt NewEmployee) fakeData(counter int) NewEmployee {
	return NewEmployee{
		Name:     "Sachin Prasad",
		Position: "Senior Software Engineer",
	}
}

// Seed runs create methods from an array of new values
func (c Core) Seed(ctx context.Context, data []NewEmployee) error {
	now := time.Now().UTC()
	for _, ns := range data {
		if _, err := c.Create(ctx, ns, now); err != nil {
			return fmt.Errorf("error seeding status: %w", err)
		}
	}

	return nil
}
