package sqlx

import (
	"fmt"
)

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
	fmt.Printf("%q\n", buf)

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

