package sqlx

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
	c Tree,
) Query {
	q.Condition = c
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

		valuers := q.Condition.values()
		for _, v := range valuers {
			vals = append(vals, any(v))
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
	case AlterColumnTypeQueryType :
		vals := []any{q.Columns[0].Default}
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

