package sqlx

import (
	"fmt"
)

type TableSchema struct {
	Name string
	Fields []TableField
}

type TableField struct {
	Name string
	Type string
	Nullable bool
	Key string
	Default string
	Extra string
}

func (db* DB)GetTableSchemas() ([]TableSchema, error) {
	var (
		ret []TableSchema
	)

	ret = []TableSchema{}

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
		"IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT " +
		"from INFORMATION_SCHEMA.COLUMNS " +
		"where TABLE_NAME=? "+
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
		)
		if nullable == "YES" {
			field.Nullable = true
		} 

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

func (db *DB)FieldToSQL(f TableField) string {
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

	ret += " " + f.Extra

	return ret
}

func (db *DB)StructToTableName(v any) string {
	return "yes"
}
