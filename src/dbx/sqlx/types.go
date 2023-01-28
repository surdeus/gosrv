package sqlx

import (
	"database/sql/driver"
	"database/sql"
	"reflect"
)

type Valuer = driver.Valuer

type NullByte = sql.NullByte
type NullInt16 = sql.NullInt16
type NullInt32 = sql.NullInt32
type NullInt = NullInt32
type NullInt64 = sql.NullInt64
type NullString = sql.NullString
type NullTime = sql.NullTime
type NullFloat64 = sql.NullFloat64
type NullFloat = NullFloat64

type Valuers []Valuer
type SqlType int

type Sqler interface {
	Sql() *TableSchema
}
type Sqlers []Sqler

type TableSchema struct {
	OldName TableName
	Name TableName
	Columns Columns
	ColMap ColumnMap
	Type reflect.Type
}

type Column struct {
	OldName ColumnName
	Name ColumnName
	Type ColumnType
	Nullable bool
	Key Key
	Default Valuer
	Extra ExtraColInfo
}
type Columns []*Column

type TableSchemas []*TableSchema

type ColumnMap map[ColumnName] *Column
type TableMap map[TableName] *TableSchema
type TableColumnMap map[TableName] ColumnMap
type TypeMap map[TableName] reflect.Type

type KeyType int
type Key struct {
	Type KeyType
}

type ExtraColInfo struct {
	AutoIncrement bool
}

// The interface type must implement to be converted to
// SQL code to be inserted for safety.
type Rawer interface {
	SqlRaw(db *Db) (Raw, error)
}

type Rawers []Rawer

type ColumnName string
type ColumnNames []ColumnName
type TableName string 
type TableNames []TableName

// Type to save raw strings for substitution
// in string queries.
type Raw string
type QueryType int
type QueryFormatFunc func(*Db, Query) (Raw, error)

type Query struct {
	Type QueryType
	Tables TableSchemas
	ColumnNames ColumnNames
	TableNames TableNames
	Columns Columns
	ColumnTypes []ColumnType
	Condition Condition
	Valuers Valuers
}

type ColumnVarType int
type ColumnType struct {
	VarType ColumnVarType
	Args []int
}

type ConditionOp int
type Condition struct {
	Op ConditionOp
	Column ColumnName
	Values Valuers
	Pair [2]*Condition
}
type Conditions []Condition

type Result struct {
	LastInsertId,
	RowsAffected int64
}
