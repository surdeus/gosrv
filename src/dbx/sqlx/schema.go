package sqlx

import (
	"strings"
	"fmt"
	"errors"
	"log"
	"github.com/surdeus/godat/src/mapx"
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

type ColumnVarType int
type ColumnType struct {
	VarType ColumnVarType
	Args RawValuers
}

type KeyType int
type Key struct {
	Type KeyType
}

type Column struct {
	OldName ColumnName
	Name ColumnName
	Type ColumnType
	Nullable bool
	Key Key
	Default RawValuer
	Extra Code
}

type Columns []Column

const (
	NotKeyType KeyType = iota
	PrimaryKeyType
	UniqueKeyType
	ForeignKeyType
)

const (
	NoColumnVarType = iota
	IntColumnVarType

	BitColumnVarType
	TinyintColumnVarType

	VarcharColumnVarType
	NvarcharColumnVarType

	CharColumnVarType
	NcharColumnVarType

	TextColumnVarType
	NtextColumnVarType

	DateColumnVarType
	TimeColumnVarType
	TimestampColumnVarType
	DatetimeColumnVarType
	YearColumnVarType

	BinaryColumnVarType
	VarbinaryColumnVarType

	ImageColumnVarType

	ClobColumnVarType
	BlobColumnVarType
	XmlColumnVarType
	JsonColumnVarType
)

var (
	MultiplePrimaryKeysErr = errors.New("multiple primary keys")
	NoPrimaryKeySpecifiedErr = errors.New("no primary key specified")
	UnknownKeyTypeErr = errors.New(
		"unknown key type",
	)
	UnknownColumnTypeErr = errors.New(
		"unknown column type",
	)
	WrongColumnTypeFormatErr = errors.New(
		"wrong column type format",
	)

	MysqlStringMapKeyType = map[string] KeyType {
		"" : NotKeyType,
		"PRI" : PrimaryKeyType,
		"UNI" : UniqueKeyType,
		"MUL" : ForeignKeyType,
	}
	MysqlKeyTypeMapString = mapx.Reverse(
		MysqlStringMapKeyType,
	)
	MysqlColumnTypeMapString = map[ColumnVarType] string {
		IntColumnVarType : "int",

		BitColumnVarType : "bit",
		TinyintColumnVarType : "tinyint",

		VarcharColumnVarType : "varchar",
		NvarcharColumnVarType : "nvarchar",

		CharColumnVarType : "char",
		NcharColumnVarType : "nchar",

		TextColumnVarType : "text",
		NtextColumnVarType : "ntext",

		DateColumnVarType : "date",
		TimeColumnVarType : "time",
		TimestampColumnVarType : "timestamp",
		DatetimeColumnVarType : "datetime",
		YearColumnVarType : "year",

		BinaryColumnVarType : "binary",
		VarbinaryColumnVarType : "varbinary",

		ImageColumnVarType : "image",

		ClobColumnVarType : "clob",
		BlobColumnVarType : "blob",
		XmlColumnVarType : "xml",
		JsonColumnVarType : "json",
	}

	MysqlStringMapColumnType = mapx.Reverse(
		MysqlColumnTypeMapString,
	)
)

func (cs Columns)Names() ColumnNames {
	ret := ColumnNames{}
	for _, v := range cs {
		ret = append(ret, v.Name)
	}

	return ret
}

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

func (db *DB)ParseColumnType(
	t string,
) (ColumnType, error) {
	ret := ColumnType{}
	t = strings.ReplaceAll(t, " ", "")
	varTypeStr, argStr, f := strings.Cut(t, "(")
	if f {
		if argStr[len(argStr)-1] != ')' {
			return ColumnType{},
				WrongColumnTypeFormatErr
		}
	}

	varTypeStr =
		strings.ToLower(varTypeStr)
	varType, ok :=
		MysqlStringMapColumnType[varTypeStr]
	if !ok {
		return ColumnType{},
			UnknownColumnTypeErr
	}

	args := RawValuers{}
	argStrs := strings.Split(argStr, ",")
	for _, v := range argStrs {
		args = append(args, RawValue(v))
	}

	ret.VarType = varType
	ret.Args = args

	return ret, nil
}

func (db *DB)GetColumnsByTableName(name TableName) (Columns, error) {
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

	var (
		cname, t, def, extra,
			key, nullable string
	)
	for rows.Next() {
		column := Column{}
		rows.Scan(
			&cname,
			&t,
			&nullable,
			&key,
			&def,
			&extra,
		)

		column.Name = ColumnName(cname)

		if nullable == "YES" {
			column.Nullable = true
		} 

		keyType, ok := MysqlStringMapKeyType[key]
		if !ok {
			return Columns{}, UnknownKeyTypeErr
		}
		column.Key.Type = keyType

		ret = append(ret, column)
	}


	return ret, nil
}

func (f Column)String() string {
	t, err := f.Type.Code()
	if err != nil {
		log.Println(err)
		return ""
	}

	def, err := f.Default.SqlRawValue()
	if err != nil {
		log.Println(err)
		return ""
	}

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
		string(t),
		f.Nullable,
		f.Key.Type,
		string(def),
		f.Extra,
	)
}

func (db *DB)ColumnToSql(f Column) (Code, error) {
	name, err := f.Name.SqlRawValue()
	if err != nil {
		return "", err
	}

	t, err := f.Type.Code()
	if err != nil {
		return "", err
	}

	ret := fmt.Sprintf(
		"%s %s",
		name,
		t,
	)

	if !f.Nullable {
		ret += " not null"
	}

	switch f.Key.Type {
	case PrimaryKeyType :
		ret += " primary key"
	default:
	}

	if string(f.Extra) != "" {
		ret += " " + string(f.Extra)
	}

	
	var (
		def RawValue
	)
	if f.Default != nil {
		def, err = f.Default.SqlRawValue()
		if err != nil {
			return "", err
		}
	}
	if def != "" {
		ret += " default " + string(def)
	}

	return Code(ret), nil
}

func (db *DB)TableCreationStringForSchema(ts TableSchema) (string, error) {
	ret := fmt.Sprintf("create table %s (\n", ts.Name)
	for i, f := range ts.Columns{
		sql, err := db.ColumnToSql(f)
		if err != nil {
			return "", err
		}
		ret += "\t" + string(sql)
		if i != len(ts.Columns) - 1 {
			ret += ",\n"
		} 
	}

	ret += "\n) ;"	

	return ret, nil
}

func (db *DB)TableCreationStringFor(v Sqler) (string, error) {
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

