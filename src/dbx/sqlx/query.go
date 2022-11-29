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
	fn, ok := queryFormatMap[q.Type]
	if !ok {
		return "", UnknownQueryTypeErr
	}
	return fn(db, q)
}

func (q Query)wType(t QueryType) Query {
	q.Type = t
	return q
}

func (q Query)Select(cn ...ColumnName) Query {
	q.ColumnNames = cn
	return q.wType(SelectQueryType)
}

func (q Query)From(table TableName) Query {
	q.TableNames = TableNames{table}
	return q
}

func (q Query)Insert(cn ...ColumnName) Query {
	q.ColumnNames = cn
	return q.wType(InsertQueryType)
}

func (q Query)Where(
	cn ColumnName,
	op ConditionOp,
	vs ...Valuer,
) Query {
	q.Conditions = Conditions{{cn, op, vs}}
	return q
}

func (q Query)And(
	cn ColumnName,
	op ConditionOp,
	vs ...Valuer,
) Query {
	q.Conditions = append(
		q.Conditions,
		Condition{cn, op, vs},
	)
	return q
}

func (q Query)CreateTable(ts *TableSchema) Query {
	q.Tables = TableSchemas{ts}
	return q.wType(CreateTableQueryType)
}

func (q Query)RenameTable(old, n TableName) Query {
	q.TableNames = TableNames{old, n}
	return q.wType(RenameTableQueryType)
}

func (q Query)RenameColumn(
	table TableName,
	old, n ColumnName,
) Query {
	q.TableNames = TableNames{table}
	q.ColumnNames = ColumnNames{old, n}
	return q.wType(RenameColumnQueryType)
}

func (q Query)AlterColumnType(
	table TableName,
	c *Column,
) Query {
	q.TableNames = TableNames{table}
	q.Columns = Columns{c}
	return q.wType(AlterColumnTypeQueryType)
}

func (q Query)Values(vs ...Valuer) Query {
	q.Valuers = vs
	return q
}

func (q Query)DropPrimaryKey(
	table TableName,
) Query {
	q.Type = DropPrimaryKeyQueryType
	q.TableNames = TableNames{table}
	return q
}

func (q Query)Into(table TableName) Query {
	q.TableNames = TableNames{table}
	return q
}

func (q Query)GetValues() []any {
	switch q.Type {
	case SelectQueryType :
		vals := []any{}
		for _, c := range q.Conditions {
			for _, v := range c.Values {
				vals = append(vals, any(v))
			}
		}
		return vals
	case InsertQueryType :
		vals := []any{}
		for _, v := range q.Valuers {
			vals = append(vals, any(v))
		}
		return vals
	case CreateTableQueryType :
		vals := []any{}
		for _, col := range q.Tables[0].Columns {
			if col.Default != nil {
				vals = append(vals, any(col.Default))
			}
		}
		return vals
	default:
		return []any{}
	}
}

func (q Query) GetColumnNames() ColumnNames {
	return q.ColumnNames
}
func (q Query) GetTableName() TableName {
	switch q.Type {
	case SelectQueryType :
		return q.TableNames[0]
	default :
		return ""
	}
}

func (q Query) WConditions(cs Conditions) Query {
	q.Conditions = cs
	return q
}

