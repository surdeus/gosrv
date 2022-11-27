package sqlx

// Simple package embedding value recievers
// for more friendly usage.

import (
	"database/sql/driver"
	"database/sql"
	"time"
)

type Valuer = driver.Valuer
type Valuers []Valuer

func Int64(n int64) sql.NullInt64 {
	return sql.NullInt64{n, true}
}

func Int(n int32) sql.NullInt32 {
	return Int32(n)
}

func Int32(n int32) sql.NullInt32 {
	return sql.NullInt32{n, true}
}

func Int16(n int16) sql.NullInt16 {
	return sql.NullInt16{n, true}
}

func Byte(n byte) sql.NullByte {
	return sql.NullByte{n, true}
}

func String(n string) sql.NullString {
	return sql.NullString{n, true}
}

func Float(n float64) sql.NullFloat64 {
	return sql.NullFloat64{n, true}
}

func Time(t time.Time) sql.NullTime {
	return sql.NullTime{t, true}
}

