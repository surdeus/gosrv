package sqlx

// Simple package embedding value recievers
// for more friendly usage.

import (
	"database/sql"
	"time"
	"strconv"
	"errors"
)

func Bool(b bool) sql.NullBool {
	return sql.NullBool{b, true}
}

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

func Float64(n float64) sql.NullFloat64 {
	return sql.NullFloat64{n, true}
}

func Float(n float64) sql.NullFloat64 {
	return Float64(n)
}

func Time(t time.Time) sql.NullTime {
	return sql.NullTime{t, true}
}

func Null() sql.NullByte {
	return sql.NullByte{0, false}
}

func StringToValuer(
	v string,
	cvt ColumnVarType,
) (Valuer, error) {
	t, ok := VarTypeMapSqlType[cvt]
	if !ok {
		return nil, UnknownColumnTypeErr
	}

	switch t {
	case BoolSqlType :
		b, err := strconv.ParseBool(v)
		return Bool(b), err
	case ByteSqlType :
		b, err := strconv.ParseInt(v, 0, 8)
		return Byte(byte(b)), err
	case Int16SqlType :
		i, err := strconv.ParseInt(v, 0, 16)
		return Int16(int16(i)), err
	case Int32SqlType :
		i, err := strconv.ParseInt(v, 0, 32)
		return Int32(int32(i)), err
	case Int64SqlType :
		i, err := strconv.ParseInt(v, 0, 64)
		return Int64(int64(i)), err
	case Float64SqlType :
		f, err := strconv.ParseFloat(v, 64)
		return Float64(f), err
	case StringSqlType :
		return String(v), nil
	default:
		return nil, errors.New("Not implemented conversion")
	}

}

