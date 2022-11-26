package sqlx

import (
	//"fmt"
)

type ColumnDiff int

const (
	NoColumnDiff ColumnDiff = iota
	NameColumnDiff
	TypeColumnDiff
	NullableColumnDiff
	KeyColumnDiff
	DefaultColumnDiff
	ExtraColumnDiff
)

func (db *DB)CompareColumns(
	c1, c2 *Column,
) (ColumnDiff, error) {

	if c1.Name != c2.Name {
		return NameColumnDiff, nil
	}

	eq, err := db.CodersEq(c1.Type, c2.Type)
	if err != nil {
		return NoColumnDiff, err
	}
	if !eq {
		return TypeColumnDiff, nil
	}

	if c1.Nullable != c2.Nullable {
		return NullableColumnDiff, nil
	}

	if !db.KeysEq(c1.Key, c2.Key) {
		return KeyColumnDiff, nil
	}

	eq, err = db.RawValuersEq(c1.Default, c2.Default)
	if err != nil {
		return NoColumnDiff, err
	}

	if !eq {
		return DefaultColumnDiff, nil
	}

	eq, err = db.CodersEq(c1.Extra, c2.Extra)
	if err != nil {
		return NoColumnDiff, err
	}
	if !eq {
		return ExtraColumnDiff, nil
	}

	return NoColumnDiff, nil
}

// Simple migration. Fuck writing much of shit in SQL directly.

func (db *DB)Migrate(sqlers []Sqler) error {
	var err error

	//curSchemas, err := db.GetTableSchemas()
	if err != nil {
		return err
	}
	newSchemas := TableSchemas{}
	for _, sqler := range sqlers {
		newSchemas = append(newSchemas, sqler.Sql())
	}

	return nil
}

/*func (db *DB)MigrateSchema(schema *TableSchema) error {
	var (
		tableName TableName
	)

	if schema.OldName != "" {
		tableName = schema.OldName
	} else {
		if schema.Name == "" {
			return WrongValuerFormatErr
		}
		tableName = schema.Name
	}

	for {
		diff, err := db.CompareColumns()
		if err != nil {
			return err
		}

		if diff == NoColumnDiff {
			break
		}

		switch diff {
		case NameColumnDiff :
		}
	}

	return nil
}

func (db *DB)MigrateRenameTable(ts *TableSchema) error {
	if ts.OldName == "" {
		return nil
	}

	
}*/

