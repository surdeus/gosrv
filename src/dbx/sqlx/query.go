package sqlx

import (
	//"fmt"
	//"strings"
	"errors"
	//"database/sql"
	//"database/sql/driver"
	//"strconv"
)


var (
	NoTablesSpecifiedErr = errors.New("no table specified")
	NoColumnsSpecifiedErr = errors.New("no columns specified")
	WrongNumOfColumnsSpecifiedErr = errors.New(
		"wrong number of columns specified")
	WrongQueryInputFormatErr = errors.New(
		"wrong query input format",
	)
	WrongValuerFormatErr = errors.New("wrong valuer format")
	UnknownQueryTypeErr = errors.New("unknown query type")
	UnknownConditionOpErr = errors.New("unknown condition operator")
	NoDBSpecifiedErr = errors.New("no database specified")
	NoSchemaSpecifiedErr = errors.New("no schema specified")

)

func Q() Query {
	return Query{}
}

func (q Query)SqlRaw(db *Db) (Raw, error) {
	fn, ok := queryFormatMap[q.typ]
	if !ok {
		return "", UnknownQueryTypeErr
	}
	return fn(db, q)
}

func (q Query)wType(t QueryType) Query {
	q.typ = t
	return q
}

func (q Query)Select(cn ...ColumnName) Query {
	q.columnNames = cn
	return q.wType(SelectQueryType)
}

func (q Query)From(table TableName) Query {
	q.tableNames = TableNames{table}
	return q
}

func (q Query)Insert() Query {
	return q.wType(InsertQueryType)
}

func (q Query)Where(c ...Condition) Query {
	q.conditions = c
	return q
}

func (q Query)And(c ...Condition) Query {
	q.conditions = append(q.conditions, c...)
	return q
}

func (q Query)CreateTable(ts *TableSchema) Query {
	q.tableSchemas = TableSchemas{ts}
	return q.wType(CreateTableQueryType)
}

func (q Query)RenameTable(old, n TableName) Query {
	q.tableNames = TableNames{old, n}
	return q.wType(RenameTableQueryType)
}

func (q Query)RenameColumn(
	table TableName,
	old, n ColumnName,
) Query {
	q.tableNames = TableNames{table}
	q.columnNames = ColumnNames{old, n}
	return q.wType(RenameColumnQueryType)
}

func (q Query)AlterColumnType(
	table TableName,
	c *Column,
) Query {
	q.tableNames = TableNames{table}
	q.columns = Columns{c}
	return q.wType(AlterColumnTypeQueryType)
}

