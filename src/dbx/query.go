package dbx

// The interfaces representns functions that
// relational database query must have based on
// vague SQL specs to implement really structured
// queries.

type RelQuerer interface {
	Create() error
	Alter() error
	Drop() error

	Select() (chan any, error)
	Insert() error
	Update() error
	Delete() error

	Grant() error
	Revoke() error
	Deny() error

	Commit() error
	Rollback() error
	Savepoint() error
}

