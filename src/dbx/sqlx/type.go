package sqlx

// Simple package embedding value recievers
// for more friendly usage.

import (
	"database/sql"
	"errors"

	//"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	case BoolSqlType:
		b, err := strconv.ParseBool(v)
		return Bool(b), err
	case ByteSqlType:
		b, err := strconv.ParseInt(v, 0, 8)
		return Byte(byte(b)), err
	case Int16SqlType:
		i, err := strconv.ParseInt(v, 0, 16)
		return Int16(int16(i)), err
	case Int32SqlType:
		i, err := strconv.ParseInt(v, 0, 32)
		return Int32(int32(i)), err
	case Int64SqlType:
		i, err := strconv.ParseInt(v, 0, 64)
		return Int64(int64(i)), err
	case Float64SqlType:
		f, err := strconv.ParseFloat(v, 64)
		return Float64(f), err
	case StringSqlType:
		v = strings.Replace(v, "''", "'", -1)
		return String(v[1 : len(v)-1]), nil
	default:
		return nil, errors.New("Not implemented conversion")
	}

}

func ValuerIsValid(v Valuer) bool {
	val := reflect.ValueOf(v)
	return val.FieldByName("Valid").
		Interface().(bool)
}

func ValueOf(v Valuer) any {
	if v == nil || !ValuerIsValid(v) {
		return nil
	}

	switch v.(type) {
	case sql.NullBool:
		vs := v.(sql.NullBool)
		return vs.Bool
	case sql.NullByte:
		vs := v.(sql.NullByte)
		return vs.Byte
	case sql.NullInt16:
		vs := v.(sql.NullInt16)
		return vs.Int16
	case sql.NullInt32:
		vs := v.(sql.NullInt32)
		return vs.Int32
	case sql.NullInt64:
		vs := v.(sql.NullInt64)
		return vs.Int64
	case sql.NullFloat64:
		vs := v.(sql.NullFloat64)
		return vs.Float64
	case sql.NullString:
		vs := v.(sql.NullString)
		return vs.String
	case sql.NullTime:
		vs := v.(sql.NullTime)
		return vs.Time
	}

	return nil
}

func ToValuer(
	v any,
) Valuer {
	switch v.(type) {
	case string:
		return String(v.(string))
	case byte:
		return Byte(v.(byte))
	case int16:
		return Int16(v.(int16))
	case int32:
		return Int32(v.(int32))
	case int64:
		return Int64(v.(int64))
	case float64:
		return Float(v.(float64))
	case time.Time:
		return Time(v.(time.Time))
	}
	return nil
}

func ValuerToString(
	v Valuer,
	format ...any,
) string {
	if !ValuerIsValid(v) {
		return "null"
	}

	switch v.(type) {
	case sql.NullBool:
		vs := v.(sql.NullBool)
		return strconv.FormatBool(vs.Bool)
	case sql.NullByte:
		vs := v.(sql.NullByte)
		return strconv.FormatInt(int64(vs.Byte), 10)
	case sql.NullInt16:
		vs := v.(sql.NullInt16)
		return strconv.FormatInt(int64(vs.Int16), 10)
	case sql.NullInt32:
		vs := v.(sql.NullInt32)
		return strconv.FormatInt(int64(vs.Int32), 10)
	case sql.NullInt64:
		vs := v.(sql.NullInt64)
		return strconv.FormatInt(vs.Int64, 10)
	case sql.NullFloat64:
		vs := v.(sql.NullFloat64)
		return strconv.FormatFloat(
			vs.Float64,
			'f', 5, 64)
	case sql.NullString:
		vs := v.(sql.NullString)
		return vs.String
	case sql.NullTime:
		vs := v.(sql.NullBool)
		return strconv.FormatBool(vs.Bool)
	default:
		return ""
	}
}

func ValuersEq(v1, v2 Valuer) bool {
	return ValueOf(v1) == ValueOf(v2)
}
