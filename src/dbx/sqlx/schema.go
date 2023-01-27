package sqlx

import (
	"strings"
	"fmt"
	"log"
	"database/sql"
	"strconv"
	"reflect"
)

const (
	NotKeyType KeyType = iota
	PrimaryKeyType
	UniqueKeyType
	ForeignKeyType
)

func (db *Db) ConstructInsertValue(
	q Query,
) (Sqler, error) {
	if q.Type != InsertQueryType {
		return nil, nil
	}

	_, err := q.SqlRaw(db)
	if err != nil {
		return nil, err
	}

	tableName := q.TableNames[0]
	columnNames := q.ColumnNames
	valuers := q.Valuers
	table, ok := db.TMap[tableName]
	if !ok {
		return nil, TableDoesNotExistErr
	}

	ret := reflect.New(table.Type).Elem()
	for i, name := range columnNames {
		f := ret.FieldByName(string(name))
		v := valuers[i]

		t1 := reflect.TypeOf(f.Interface())
		t2 := reflect.TypeOf(v)

		if !t1.AssignableTo(t2) {
			return nil, NotAssignableErr
		}

		f.Set(reflect.ValueOf(valuers[i]))
	}

	return ret.Interface().(Sqler), nil
}

func TableBySqler(sqler Sqler) *TableSchema {
	ret := sqler.Sql()

	if ret.Type == nil {
		ret.Type = reflect.TypeOf(sqler)
	}

	ret.ColMap = ret.Columns.ColumnMap()

	return ret
}

func (sqlers Sqlers)Tables() TableSchemas {
	ret := TableSchemas{}
	for _, sqler := range sqlers {
		t := TableBySqler(sqler)
		ret = append(ret, t)
	}

	return ret
}

func (tables TableSchemas)TypeMap() TypeMap {
	ret := make(TypeMap)
	for _, table := range tables {
		ret[table.Name] = table.Type
	}

	return ret
}

func (tables TableSchemas)TableMap() TableMap {
	ret := make(TableMap)
	for _, table := range tables {
		ret[table.Name] = table
	}

	return ret
}

func (tables TableSchemas)TableColumnMap() TableColumnMap {
	ret := make(TableColumnMap)
	for _, table := range tables {
		ret[table.Name] =
			table.Columns.ColumnMap()
	}

	return ret
}

func (cols Columns)ColumnMap() ColumnMap {
	ret := make(ColumnMap)
	for _, col := range cols {
		ret[col.Name] = col
	}

	return ret
}

func (cs Columns)Names() ColumnNames {
	ret := ColumnNames{}
	for _, v := range cs {
		ret = append(ret, v.Name)
	}

	return ret
}

func NotKey() Key {
	return Key {
		Type : NotKeyType,
	}
}

func PrimaryKey() Key {
	return Key{Type: PrimaryKeyType}
}

func (schema *TableSchema)PrimaryKeyColumn() (int, *Column, error) {
	var (
		ret, i, n int
		 f *Column
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

	return ret, schema.Columns[ret], nil
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
			return i, schemas[i]
		}
	}

	return -1, nil
}

func (ts TableSchema)FindColumn(
	name ColumnName,
) (int, *Column) {
	for i, _ := range ts.Columns {
		if ts.Columns[i].Name == name {
			return i, ts.Columns[i]
		}
	}

	return -1, nil
}

func (db *Db)GetTableNames(
) (TableNames, error) {
	var (
		ret TableNames
	)

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
		name := TableName("")
		rows.Scan(&name)
		ret = append(ret, name)
	}

	return ret, nil
}

func (db *Db)GetTableSchema(
	name TableName,
) (*TableSchema, error) {
	var err error
	exists, err := db.TableExists(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, TableDoesNotExistErr
	}

	ret := &TableSchema{}
	ret.Name = name
	ret.Columns, err = db.GetColumnsByTableName(name)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *Db)ColumnFromRaws(
	cname, t, nullable,
	key, extra string,
	def sql.NullString,
) (*Column, error) {
	var err error
	column := new(Column)

	column.Name = ColumnName(cname)
	column.Type, err = db.ParseColumnType(t)
	if err != nil {
		return nil, err
	}

	if nullable == "YES" {
		column.Nullable = true
	} 

	keyType, ok := MysqlStringMapKeyType[key]
	if !ok {
		return nil, UnknownKeyTypeErr
	}
	column.Key.Type = keyType

	if !def.Valid {
		column.Default = nil
	} else {
		valuer, err := StringToValuer(
			def.String, column.Type.VarType)
		if err != nil {
			return nil, err
		}
		column.Default = valuer
	}
	column.Extra = Raw(extra)

	return column, nil
}

func (db *Db)GetColumnSchema(
	table TableName,
	colName ColumnName,
) (*Column, error) {
	exists, err := db.ColumnExists(table, colName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ColumnDoesNotExistErr
	}

	rows, err := db.Query(
		"select "+
		"COLUMN_NAME, COLUMN_TYPE, " +
		"IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT, EXTRA " +
		"from INFORMATION_SCHEMA.COLUMNS " +
		"where TABLE_NAME in (?) and COLUMN_NAME in (?)"+
		"",
		table,
		colName,
	)
	if err != nil {
		return nil, err
	}

	var (
		cname, t, extra,
			key, nullable string
		def sql.NullString
	)

	/*if !rows.Next() {
		return nil, ColumnDoesNotExistErr
	}*/
	rows.Next()
	err = rows.Scan(
		&cname,
		&t,
		&nullable,
		&key,
		&def,
		&extra,
	)
	if err != nil {
		return nil, err
	}

	return db.ColumnFromRaws(
		cname, t, nullable,
		key, extra,
		def,
	)
}


func (db* Db)GetTableSchemas() (TableSchemas, error) {
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
		s := &TableSchema{}

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

func (db *Db)TableExists(name TableName) (bool, error) {
	ret := false
	raw, err := name.SqlRaw(db)
	if err != nil {
		return false, err
	}

	rows, err := db.Query(fmt.Sprintf(
		"select * from %s ;",
		raw,
	))

	if err == nil {
		defer rows.Close()
		ret = true
	}

	return ret, nil
}

func (db *Db)RenameTable(old, n TableName) error {
	_, _, err := db.Do(
		Q().RenameTable(old, n),
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db)RenameColumn(
	table TableName,
	o, n ColumnName,
) error {
	_, _, err := db.Do(
		Q().RenameColumn(table, o, n),
	)
	return err
}

func (db *Db)ParseColumnType(
	t string,
) (ColumnType, error) {
	ret := ColumnType{}
	t = strings.ReplaceAll(t, " ", "")
	varTypeStr, argStr, f := strings.Cut(t, "(")
	if f {
		if argStr[len(argStr)-1] != ')' {
			return ColumnType{},
				fmt.Errorf("ParseColumn: %s: %w",
					t, WrongColumnTypeFormatErr )
		}
	}

	varTypeStr =
		strings.ToLower(varTypeStr)
	varType, ok :=
		MysqlStringMapColumnType[varTypeStr]
	if !ok {
		return ColumnType{},
			fmt.Errorf("ParseColumn: %s: %w",
					t, UnknownColumnTypeErr)
	}

	if f {
		argStr = argStr[:len(argStr)-1]
	}
	args := []int{}
	argStrs := strings.Split(
		argStr,
		",",
	)
	for _, v := range argStrs {
		i, err := strconv.Atoi(v)
		if err != nil {
			return ColumnType{}, err
		}
		args = append(args, i)
	}

	ret.VarType = varType
	ret.Args = args

	return ret, nil
}

func (db *Db)GetColumnsByTableName(name TableName) (Columns, error) {
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
		cname, t, extra,
			key, nullable string
		def sql.NullString
	)
	for rows.Next() {
		column := &Column{}
		rows.Scan(
			&cname,
			&t,
			&nullable,
			&key,
			&def,
			&extra,
		)

		column, err := db.ColumnFromRaws(
			cname, t, nullable,
			key, extra,
			def,
		)
		if err != nil {
			return nil, err
		}

		ret = append(ret, column)
	}


	return ret, nil
}

func (f Column)String() string {
	t, err := f.Type.SqlRaw(nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	def, err := f.Default.Value()

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
		def,
		f.Extra,
	)
}

func (db *Db)ColumnToAlterSql(
	c *Column,
) (Raw, error) {
	buf := *c
	buf.Key = NotKey()
	return db.ColumnToSql(&buf)
}

func (db *Db)ColumnToSql(f *Column) (Raw, error) {
	ret, err := db.Rprintf(
		"%s %s",
		f.Name,
		f.Type,
	)
	if err != nil {
		return "", err
	}

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
	
	if f.Default != nil {
		ret += " default ?"
	}

	return Raw(ret), nil
}

func (ts *TableSchema)SqlRaw(db *Db) (Raw, error) {
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

	return Raw(ret), nil
}

func (db *Db)TableCreationStringFor(v Sqler) (string, error) {
	c, err := v.Sql().SqlRaw(db)
	return string(c), err
}

func (db *Db)CreateTable(v Sqler) error {
	_, _, err := db.Do(
		Q().CreateTable(v.Sql()),
	)
	return err
}

func (db *Db)CreateTableBySchema(ts *TableSchema) error {
	return nil
}

func (db *Db)AlterAddColumn(
	tn TableName, c *Column,
) error {
	var err error
	table, err := tn.SqlRaw(db)
	if err != nil {
		return err
	}
	t, err := db.ColumnToAlterSql(c)
	if err != nil {
		return err
	}

	rows, err := db.Query(
	fmt.Sprintf(
		"alter table %s add %s;",
		table,
		t,
	))
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (db *Db)AlterColumnType(
	table TableName,
	c *Column,
) error {
	_, _, err := db.Do(
		Q().AlterColumnType(table, c),
	)

	return err
}

func (db *Db)ColumnExists(
	table TableName,
	column ColumnName,
) (bool, error) {
	traw, err := table.SqlRaw(db)
	if err != nil {
		return false, err
	}

	craw, err := column.SqlRaw(db)
	if err != nil {
		return false, err
	}

	rows, err := db.Query(fmt.Sprintf(
		"select %s from %s limit 1 ;",
		craw,
		traw,
	))
	if err == nil {
		rows.Close()
		return true, nil
	}

	return false, nil
}

func (db *Db)DropPrimaryKey(name TableName) error {
	q := Q().DropPrimaryKey(name)
	_, _, err := db.Do(q)

	return err
}

func (db *Db)KeysEq(k1, k2 Key) bool {
	return k1.Type == k2.Type
}

