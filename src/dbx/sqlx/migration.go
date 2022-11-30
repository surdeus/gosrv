package sqlx

import (
	"errors"
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

var (
	OldAndNewTablesExistErr = errors.New("old and new tables exist")
)

/*func (db *Db)CompareColumns(
	c1, c2 *Column,
) (ColumnDiff, error) {

	if c1.Name != c2.Name {
		return NameColumnDiff, nil
	}

	eq, err := db.RawersEq(c1.Type, c2.Type)
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

	eq, err = db.RawersEq(c1.Default, c2.Default)
	if err != nil {
		return NoColumnDiff, err
	}

	if !eq {
		return DefaultColumnDiff, nil
	}

	eq, err = db.RawersEq(c1.Extra, c2.Extra)
	if err != nil {
		return NoColumnDiff, err
	}
	if !eq {
		return ExtraColumnDiff, nil
	}

	return NoColumnDiff, nil
}

*/

// Simple migration. Fuck writing much of shit in SQL directly.

func (db *Db)Migrate() error {
	var err error

	if err != nil {
		return err
	}
	schemas := db.Tables

	for _, schema := range schemas {
		err = db.MigrateSchema(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Db)MigrateSchema(
	schema *TableSchema,
) error {
	err := db.MigrateRenameTable(schema)
	if err != nil &&
			err != TableAlreadyExistsErr {

		if err == TableDoesNotExistErr {
			err := db.CreateTableBySchema(schema)
			return err
		}

		return err
	}

	for _, c := range schema.Columns {
		var err error
		exists, err := db.ColumnExists(schema.Name, c.Name)
		if err != nil {
			return err
		}
		if !exists {
			err = db.AlterAddColumn(schema.Name, c)
			if err != nil {
				return err
			}
			continue
		}
		err = db.
			MigrateRenameColumn(schema.Name, c)
		if err != nil &&
			err != ColumnAlreadyExistsErr {
			return err
		}

		err = db.MigrateAlterColumnType(schema.Name, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Db)MigrateAlterColumnType(
	tableName TableName,
	newCol *Column,
) (error) {
	curCol, err := db.GetColumnSchema(
		tableName, newCol.Name,
	)
	if err != nil {
		return err
	}

	curSql, err := db.ColumnToAlterSql(curCol)
	if err != nil {
		return err
	}

	newSql, err := db.ColumnToAlterSql(newCol)
	if err != nil {
		return err
	}

	if curSql != newSql ||
			curCol.Default != newCol.Default {
		err = db.AlterColumnType(
			tableName, newCol,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Db)MigrateRenameColumn(
	tableName TableName,
	column *Column,
) (error) {
	if column.OldName == "" {
		return nil
	}

	existsOld, err := db.ColumnExists(tableName,
		column.OldName)
	if err != nil {
		return err
	}

	existsNew, err := db.ColumnExists(tableName,
		column.Name)
	if err != nil {
		return err
	}

	if existsNew {
		return ColumnAlreadyExistsErr
	}

	if !existsOld {
		return ColumnDoesNotExistErr
	}

	err = db.RenameColumn(
		tableName,
		column.OldName,
		column.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db)MigrateRenameTable(
	ts *TableSchema,
) (error) {
	if ts.OldName == "" {
		return nil
	}

	existsOld, err := db.TableExists(ts.OldName)
	if err != nil {
		return err
	}

	existsNew, err := db.TableExists(ts.Name)
	if err != nil {
		return err
	}

	if existsNew {
		return TableAlreadyExistsErr
	}

	if !existsOld {
		return TableDoesNotExistErr
	}

	err = db.RenameTable(ts.OldName, ts.Name)
	if err != nil {
		return err
	}

	return nil
}

