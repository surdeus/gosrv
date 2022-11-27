package sqlx

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
}

const (
	NoQueryType QueryType = iota
	SelectQueryType
	InsertQueryType
	DeleteQueryType
	RenameTableQueryType
	RenameColumnQueryType
	CreateTableQueryType
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
	}
)

func insertQuery(
	db *Db,
	q Query,
) (Raw, error) {
	return Raw(""), nil
}

func selectQuery(
	db *Db,
	q Query,
) (Raw, error) {
	if len(q.TableNames) != 1 {
		return "", NoTablesSpecifiedErr
	}

	if len(q.ColumnNames) < 1 {
		return "", NoColumnsSpecifiedErr
	}

	if len(q.Conditions) > 1 {
		return "", WrongQueryInputFormatErr
	} else if len(q.Conditions) >= 1 {
		return db.Rprintf(
			"select %s from %s where %s ;",
			q.ColumnNames[0],
			q.TableNames[0],
			q.Conditions,
		)
	} else {
		return db.Rprintf(
			"select %s from %s ;",
			q.ColumnNames[0],
			q.TableNames[0],
		)
	}
	if err != nil {
		return "", err
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
	if len(q.TableSchemas) != 1 {
		return "", NoSchemaSpecifiedErr
	}

	buf, err := db.
		TableCreationStringForSchema(q.TableSchemas[0])
	if err != nil {
		return "", err
	}

	ret = Raw(buf)
	return ret, err
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
