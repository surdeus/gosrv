package ormx

import (
	"github.com/mojosa-software/gosrv/src/dbx/sqlx"
)

// The type represents way to generalize structure type.
type Object struct {
	V any
}

func (o Object) Sql() *sqlx.TableSchema {
	cols := []*sqlx.Column {
	}
}

