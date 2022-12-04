package sqlx

import (
	"database/sql/driver"
	"reflect"
)

type Valuer = driver.Valuer
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
	Type reflect.Type
}

type Column struct {
	OldName ColumnName
	Name ColumnName
	Type ColumnType
	Nullable bool
	Key Key
	Default Valuer
	Extra Raw
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
	Conditions Conditions
	Valuers Valuers
}

type ColumnVarType int
type ColumnType struct {
	VarType ColumnVarType
	Args []int
}

type ConditionOp int
type Condition struct {
	Column ColumnName
	Op ConditionOp
	Values Valuers
}
type Conditions []Condition

type Result struct {
	LastInsertId,
	RowsAffected int64
}
