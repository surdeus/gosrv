package sqlx

import (
	"fmt"
	"errors"
)

type Sqler interface {
	Sql() TableSchema
}

type TableSchema struct {
	OldName string
	Name string
	Fields TableFields
}

type TableSchemas []TableSchema

type TableField struct {
	OldName string
	Name string
	Type string
	Nullable bool
	Key string
	Default string
	Extra string
}

type TableFields []TableField

var (
	MultiplePrimaryKeysErr = errors.New("multiple primary keys")
)

func (schema TableSchema)PrimaryKeyFieldId() (int, error) {
	var (
		ret, i, n int
		 f TableField
	)
	for i, f = range schema.Fields {
		if f.IsPrimaryKey() {
			n++
			ret = i
			if n > 1 {
				return -1, MultiplePrimaryKeysErr
			}
		}
	}

	if n != 1 {
		return -1, nil
	}

	return ret, nil
}

func (f TableField)IsPrimaryKey() bool {
	return f.Key == "PRI"
}

func (schemas TableSchemas)FindSchema(name string) int {
	for i, _ := range schemas {
		if schemas[i].Name == name {
			return i
		}
	}

	return -1
}

func (ts TableSchema)FindField(name string) int {
	for i, _ := range ts.Fields {
		if ts.Fields[i].Name == name {
			return i
		}
	}

	return -1
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

		s.Fields, err = db.GetFieldsByTableName(s.Name)
		if err != nil {
			return nil, err
		}

		ret = append(ret, s)
	}

	return ret, nil
}

func (db *DB)TableExists(name string) bool {
	ret := false
	rows, err := db.Query(fmt.Sprintf("select * from %s ;", name))
	if err == nil {
		defer rows.Close()
		ret = true
	}

	return ret
}

func (db *DB)GetFieldsByTableName(name string) ([]TableField, error) {
	var (
		nullable string
	)
	ret := []TableField{}
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

	for rows.Next() {
		field := TableField{}
		rows.Scan(
			&field.Name,
			&field.Type,
			&nullable,
			&field.Key,
			&field.Default,
			&field.Extra,
		)
		if nullable == "YES" {
			field.Nullable = true
		} 

		fmt.Println(field)

		ret = append(ret, field)
	}


	return ret, nil
}

func (f TableField)String() string {
	return fmt.Sprintf(
		"{\n" +
		"\tName: \"%s\",\n" +
		"\tType: \"%s\",\n" +
		"\tNullable: %t,\n"+
		"\tKey: \"%s\",\n"+
		"\tDefault: %s,\n"+
		"\tExtra: \"%s\",\n"+
		"}",
		f.Name,
		f.Type,
		f.Nullable,
		f.Key,
		f.Default,
		f.Extra,
	)
}

func (db *DB)FieldToSql(f TableField) string {
	ret := fmt.Sprintf(
		"%s %s",
		f.Name,
		f.Type,
	)

	if !f.Nullable {
		ret += " not null"
	}

	switch f.Key {
	case "PRI" :
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
	for i, f := range ts.Fields {
		ret += "\t" + db.FieldToSql(f)
		if i != len(ts.Fields) - 1 {
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

func (db *DB)FieldExists(table, field string) bool {
	rows, err := db.Query(fmt.Sprintf("select %s from %s limit 1 ;", field, table))
	if err == nil {
		rows.Close()
		return true
	}

	return false
}


