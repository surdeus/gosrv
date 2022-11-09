package sqls

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

func (db *DB)GetFieldsByTableName(name string) ([]TableField, error) {
	ret := []TableField{}
	rows, err := db.Query(
		"select "+
		"COLUMN_NAME, COLUMN_TYPE " +
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
		)
		ret = append(ret, field)
	}


	return ret, nil
}

