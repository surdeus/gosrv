package sqlx

import (
	"fmt"
	"errors"
)

type Sqler interface {
	Sql() TableSchema
}

type TableSchema struct {
	OldName TableName
	Name TableName
	Columns Columns
}

type TableSchemas []TableSchema

type ColumnType int
type KeyType int
type Key struct {
	Type KeyType
}

type Column struct {
	OldName ColumnName
	Name ColumnName
	Type string
	TypeArgs []RawValuer
	Nullable bool
	Key Key
	Default string
	Extra string
}

type Columns []Column

const (
	NotKeyType KeyType = iota
	UniqueKeyType
	ForeignKeyType
	PrimaryKeyType
)

const (
	IntColumnType = iota

	BitColumnType
	TinyintColumnType

	VarcharColumnType
	NvarcharColumnType

	CharColumnType
	NcharColumnType

	TextColumnType
	NtextColumnType

	DateColumnType
	TimeColumnType
	TimestampColumnType
	DatetimeColumnType
	YearColumnType

	BinaryColumnType
	VarbinaryColumnType

	ImageColumnType

	ClobColumnType
	BlobColumnType
	XmlColumnType
	JsonColumnType
)

var (
	MultiplePrimaryKeysErr = errors.New("multiple primary keys")
	NoPrimaryKeySpecifiedErr = errors.New("no primary key specified")
	UnknownKeyTypeErr = errors.New(
		"unknown key type",
	)

	MysqlKeyTypeStringMap = map[string] KeyType {
		"" : NotKeyType,
		"PRI" : PrimaryKeyType,
		"UNI" : UniqueKeyType,
		"MUL" : ForeignKeyType,
	}
)

func PrimaryKey() Key {
	return Key{Type: PrimaryKeyType}
}

func (schema *TableSchema)PrimaryKeyColumn() (int, *Column, error) {
	var (
		ret, i, n int
		 f Column
	)
	for i, f = range schema.Columns {
		if f.IsPrimaryKey() {
			n++
			ret = i
			if n > 1 {
				return -1, nil, MultiplePrimaryKeysErr
			}
		}
	}

	if n != 1 {
		return -1, nil, NoPrimaryKeySpecifiedErr
	}

	return ret, &schema.Columns[ret], nil
}

func (f *Column)IsPrimaryKey() bool {
	return f.Key.Type == PrimaryKeyType
}

func (f *Column)IsNotKey() bool {
	return f.Key.Type == NotKeyType
}

func (schemas TableSchemas)FindSchema(
	name TableName,
) (int, *TableSchema) {
	for i, _ := range schemas {
		if schemas[i].Name == name {
			return i, &schemas[i]
		}
	}

	return -1, nil
}

func (ts TableSchema)FindColumn(
	name ColumnName,
) (int, *Column) {
	for i, _ := range ts.Columns {
		if ts.Columns[i].Name == name {
			return i, &(ts.Columns[i])
		}
	}

	return -1, nil
}

func (db* DB)GetTableSchemas() (TableSchemas, error) {
	var (
		ret TableSchemas
	)

	ret = TableSchemas{}

	rows, err := db.Query(
		"select " +
		"TABLE_NAME " +
		"from INFORMATION_SCHEMA.TABLES " +
		"where TABLE_SCHEMA = database() " +
		"",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := TableSchema{}

		rows.Scan(
			&s.Name,
		)

		s.Columns, err = db.GetColumnsByTableName(s.Name)
		if err != nil {
			return nil, err
		}

		ret = append(ret, s)
	}

	return ret, nil
}

func (db *DB)TableExists(name TableName) bool {
	ret := false
	rows, err := db.Query(fmt.Sprintf("select * from %s ;", name))
	if err == nil {
		defer rows.Close()
		ret = true
	}

	return ret
}

func (db *DB)GetColumnsByTableName(name TableName) (Columns, error) {
	var (
		nullable string
	)
	ret := Columns{}
	rows, err := db.Query(
		"select "+
		"COLUMN_NAME, COLUMN_TYPE, " +
		"IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT, EXTRA " +
		"from INFORMATION_SCHEMA.COLUMNS " +
		"where TABLE_NAME in (?) "+
		"",
		name,
	)
	if err != nil {
		return nil, err
	}

	var key string
	for rows.Next() {
		column := Column{}
		rows.Scan(
			&column.Name,
			&column.Type,
			&nullable,
			&key,
			&column.Default,
			&column.Extra,
		)
		if nullable == "YES" {
			column.Nullable = true
		} 
		keyType, ok := MysqlKeyTypeStringMap[key]
		if !ok {
			return Columns{}, UnknownKeyTypeErr
		}
		column.Key.Type = keyType

		fmt.Println(column)

		ret = append(ret, column)
	}


	return ret, nil
}

func (f Column)String() string {
	return fmt.Sprintf(
		"{\n" +
		"\tName: \"%s\",\n" +
		"\tType: \"%s\",\n" +
		"\tNullable: %t,\n"+
		"\tKey: %d,\n"+
		"\tDefault: %s,\n"+
		"\tExtra: \"%s\",\n"+
		"}",
		f.Name,
		f.Type,
		f.Nullable,
		f.Key.Type,
		f.Default,
		f.Extra,
	)
}

func (db *DB)ColumnToSql(f Column) string {
	ret := fmt.Sprintf(
		"%s %s",
		f.Name,
		f.Type,
	)

	if !f.Nullable {
		ret += " not null"
	}

	switch f.Key.Type {
	case PrimaryKeyType :
		ret += " primary key"
	default:
	}

	if f.Extra != "" {
		ret += " " + f.Extra
	}

	if f.Default != "" {
		ret += " default " + f.Default
	}

	return ret
}

func (db *DB)TableCreationStringForSchema(ts TableSchema) string {
	ret := fmt.Sprintf("create table %s (\n", ts.Name)
	for i, f := range ts.Columns{
		ret += "\t" + db.ColumnToSql(f)
		if i != len(ts.Columns) - 1 {
			ret += ",\n"
		} 
	}

	ret += "\n) ;"	

	return ret
}

func (db *DB)TableCreationStringFor(v Sqler) string {
	return db.TableCreationStringForSchema(v.Sql())
}

func (db *DB)CreateTableBy(v Sqler) error {
	return db.CreateTableBySchema(v.Sql())
}

func (db *DB)CreateTableBySchema(ts TableSchema) error {
	_, err := db.Query(db.TableCreationStringForSchema(ts))
	return err
}

func (db *DB)ColumnExists(
	table TableName,
	column ColumnName,
) bool {
	rows, err := db.Query(fmt.Sprintf("select %s from %s limit 1 ;", column, table))
	if err == nil {
		rows.Close()
		return true
	}

	return false
}

func (db *DB)DropTablePrimaryKey(name TableName) error {
	_, err := db.Exec(fmt.Sprintf(
		"alter table %s drop primary key ;",
		name,
	))

	if err != nil {
		return err
	}

	return nil
}

