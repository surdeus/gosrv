package dbx

import (
	"database/sql"
)

// The interface represents functions that value
// must have to be used in relational database interactions
// as concrete value, not column name or table name etc.
type RelValuer interface {
	Value() any
	IsValid() bool
}

