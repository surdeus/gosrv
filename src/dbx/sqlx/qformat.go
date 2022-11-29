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
	typ QueryType
	tableSchemas TableSchemas
	columnNames ColumnNames
	tableNames TableNames
	columns Columns
	columnTypes []ColumnType
	conditions Conditions
	values Valuers
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
	if len(q.tableNames) != 1 ||
			len(q.columnNames) < 1 ||
			len(q.columnNames) != len(q.values) {
		return "", WrongQueryInputFormatErr
	}

	r, err := db.Rprintf(
		"insert into %s (%s) values %s ;",
		q.tableNames[0],
		q.columnNames,
		db.TupleBuf(q.values),
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
	if len(q.tableNames) != 1 {
		return "", NoTablesSpecifiedErr
	}

	if len(q.columnNames) < 1 {
		return "", NoColumnsSpecifiedErr
	}

	if len(q.conditions) >= 1 {
		return db.Rprintf(
			"select %s from %s where %s ;",
			q.columnNames,
			q.tableNames[0],
			q.conditions,
		)
	} else {
		return db.Rprintf(
			"select %s from %s ;",
			q.columnNames,
			q.tableNames[0],
		)
	}
}

func renameColumn(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.columnNames) != 2 ||
			len(q.tableNames) != 1 {
		return "", WrongNumOfColumnsSpecifiedErr
	}

	return db.Rprintf(
		"alter table %s rename column %s to %s ;",
		q.tableNames[0],
		q.columnNames[0],
		q.columnNames[1],
	)
}

func renameTable(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.tableNames) != 2 {
		return "", NoTablesSpecifiedErr
	}

	return db.Rprintf(
		"alter table %s rename %s ;",
		q.tableNames[0],
		q.tableNames[1],
	)
}

func createTable(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.tableSchemas) != 1 {
		return "", NoSchemaSpecifiedErr
	}

	buf, err := q.tableSchemas[0].SqlRaw(db)
	if err != nil {
		return "", err
	}

	return Raw(buf), err
}


func alterColumnType(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.tableNames) != 1 ||
		len(q.columns) != 1 {
		return "",
			WrongQueryInputFormatErr
	}

	rcode, err := db.ColumnToAlterSql(
		q.columns[0],
	)
	if err != nil {
		return "", err
	}

	buf := fmt.Sprintf(
		"alter table %s modify %s ;",
		q.tableNames[0],
		rcode,
	)

	return Raw(buf), nil
}

func dropPrimaryKey(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.tableNames) != 1 {
		return "", WrongQueryInputFormatErr
	}
	r, err := db.Rprintf(
		"alter table %s drop primary key ;",
		q.tableNames[0],
	)

	return r, err
}

