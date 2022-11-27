package sqlx

import (
	"fmt"
	"strings"
	"errors"
	"database/sql"
	"database/sql/driver"
	"strconv"
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

func (q Query)wDatabase(db *Db) Query {
	q.DB = db
	return q
}

func (q Query)wType(t QueryType) Query {
	q.typ = t
	return q
}

func (q Query)Select(cn ...ColumnName) Query {
	q.ColumnNames 
	return q.wType(SelectQueryType)
}

func (q Query)Insert() Query {
	return q.wType(InsertQueryType)
}

func (q Query)Where(c ...Condition) Query {
	q.Conditions = c
	return q
}

func (q Query)And(c ...Condition) Query {
	q.Conditions = append(q.Conditions, c...)
	return q
}

func (q Query)CreateTable() Query {
	return q.wType(CreateTableQueryType)
}

func (q Query)RenameTable() Query {
	return q.wType(RenameTableQueryType)
}

func (q Query)RenameColumn() Query {
	return q.wType(RenameColumnQueryType)
}

func (q Query)AlterColumnType() Query {
	return q.wType(AlterColumnTypeQueryType)
}

