package dbx

// The package represents general interfaces
// for interaction with different kinds
// of databases through structured queries
// sent by Golang gob format to prevent exceed
// parsing and security issues. Operations
// and queries are based on SQL since it is
// well described already.

// The interface represents relational database
// interaction functions.
type RelDatabaser interface {
	Do() (chan any, error)
	CanDo() bool
}

