package db

import (
	"time"
)

// Employee represent the structure we need for moving data
// between the app and the database.
type Employee struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Position  string     `db:"position"`
	CreatedOn time.Time  `db:"created_on"`
	UpdatedOn time.Time  `db:"updated_on"`
	DeletedOn *time.Time `db:"deleted_on"`
}
