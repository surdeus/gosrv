package sqlx

import (
	"fmt"
)

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

const (
	NoQueryType QueryType = iota
	SelectQueryType
	InsertQueryType
	DeleteQueryType
	RenameTableQueryType
	RenameColumnQueryType
	CreateTableQueryType
	DropPrimaryKeyQueryType
	AlterColumnTypeQueryType
	ModifyQueryType
)

var (
	queryFormatMap = map[QueryType] QueryFormatFunc {
		SelectQueryType : selectQuery,
		InsertQueryType : insertQuery,
		RenameTableQueryType : renameTable,
		RenameColumnQueryType : renameColumn,
		CreateTableQueryType : createTable,
		AlterColumnTypeQueryType : alterColumnType,
		DropPrimaryKeyQueryType : dropPrimaryKey,
	}
)

func insertQuery(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.TableNames) != 1 ||
			len(q.ColumnNames) < 1 ||
			len(q.ColumnNames) != len(q.Valuers) {
		return "", WrongQueryInputFormatErr
	}

	r, err := db.Rprintf(
		"insert into %s (%s) values %s ;",
		q.TableNames[0],
		q.ColumnNames,
		db.TupleBuf(q.Valuers),
	)
	if err != nil {
		return "", err
	}

	return r, nil
}

func selectQuery(
	db *Db,
	q Query,
) (Raw, error) {

	if len(q.Conditions) >= 1 {
		return db.Rprintf(
			"select %s from %s where %s ;",
			q.ColumnNames,
			q.TableNames[0],
			q.Conditions,
		)
	} else {
		return db.Rprintf(
			"select %s from %s ;",
			q.ColumnNames,
			q.TableNames[0],
		)
	}
}

func renameColumn(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.ColumnNames) != 2 ||
			len(q.TableNames) != 1 {
		return "", WrongNumOfColumnsSpecifiedErr
	}

	return db.Rprintf(
		"alter table %s rename column %s to %s ;",
		q.TableNames[0],
		q.ColumnNames[0],
		q.ColumnNames[1],
	)
}

func renameTable(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.TableNames) != 2 {
		return "", NoTablesSpecifiedErr
	}

	return db.Rprintf(
		"alter table %s rename %s ;",
		q.TableNames[0],
		q.TableNames[1],
	)
}

func createTable(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.Tables) != 1 {
		return "", NoSchemaSpecifiedErr
	}

	buf, err := q.Tables[0].SqlRaw(db)
	if err != nil {
		return "", err
	}

	return Raw(buf), err
}


func alterColumnType(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.TableNames) != 1 ||
		len(q.Columns) != 1 {
		return "",
			WrongQueryInputFormatErr
	}

	rcode, err := db.ColumnToAlterSql(
		q.Columns[0],
	)
	if err != nil {
		return "", err
	}

	buf := fmt.Sprintf(
		"alter table %s modify %s ;",
		q.TableNames[0],
		rcode,
	)

	return Raw(buf), nil
}

func dropPrimaryKey(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.TableNames) != 1 {
		return "", WrongQueryInputFormatErr
	}
	r, err := db.Rprintf(
		"alter table %s drop primary key ;",
		q.TableNames[0],
	)

	return r, err
}

