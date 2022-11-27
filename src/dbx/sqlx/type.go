package sqlx

// Simple package embedding value recievers
// for more friendly usage.

import (
	"database/sql/driver"
	"database/sql"
)

type Valuer = driver.Valuer
type Valuers []Valuer

func Int(n int32) sql.NullInt32 {
	return sql.NullInt32{n, true}
}

func String(n string) sql.NullString {
	return sql.NullString{n, true}
}

func Float(n float64) sql.NullFloat64 {
	return sql.NullFloat64{n, true}
}

