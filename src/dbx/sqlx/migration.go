package sqlx

import (
	"fmt"
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

	curSchemas, err := db.GetTableSchemas()
	if err != nil {
		return err
	}
	newSchemas := TableSchemas{}
	for _, sqler := range sqlers {
		newSchemas = append(newSchemas, sqler.Sql())
	}

	for _, newSchema := range newSchemas {
		_, curSchema := curSchemas.FindSchema(newSchema.Name)
		for _, newColumn := range newSchema.Columns {
			i, curColumn :=
				curSchema.FindColumn(newColumn.Name)
			fmt.Println("i:", i)
			fmt.Println("names:", newColumn.Name, curColumn.Name)
			diff, err :=
				db.CompareColumns(&newColumn, curColumn)
			fmt.Println("diff: ", diff, err)
		}
	}

	return nil

	// First we should rename existing tables and create not existing ones.
	for _, schema := range newSchemas {
		// Rename.
		if schema.OldName != TableName("") &&  db.TableExists(schema.OldName) {
			q := db.Q().AlterTableRename().
				WithFrom(schema.OldName).
				WithTo(schema.Name)
			/*_, err = db.Query(fmt.Sprintf(
				"alter table %s rename %s ;",
				schema.OldName,
				schema.Name,
			))*/
			_, err = q.Do()
			if err != nil {
				return err
			}

			// Fit changes to the current schema representation.
			_, curSchema := curSchemas.FindSchema(schema.OldName)
			curSchema.Name = schema.Name

			continue
		}

		// Create.
		db.CreateTableBySchema(schema)
		curSchemas = append(curSchemas, schema)
	}

	// Then we modify existing and create not existing columns.
	for _, schema := range newSchemas {
		_, curSchema := curSchemas.FindSchema(schema.Name)
		for _, column := range schema.Columns {

			if column.OldName != ColumnName("") &&
					db.ColumnExists(schema.Name, column.OldName) {

				// Rename.
				_, curColumn := curSchema.FindColumn(column.OldName)
				//curField := &(curSchemas[idx].Fields[curFieldIdx])

				_, err = db.Query(fmt.Sprintf(
					"alter table %s rename column %s to %s ;",
					schema.Name,
					column.OldName,
					column.Name,
				))
				if err != nil {
					return err
				}

				curColumn.Name = column.Name
			} else if !db.ColumnExists(schema.Name, column.Name) {
				// Create.
				sql, err := db.ColumnToSql(column)
				if err != nil {
					return err
				}
				_, err = db.Query(fmt.Sprintf(
					"alter table %s add %s",
					schema.Name,
					sql,
				))
				if err != nil {
					return err
				}
			}

			_, curColumn := curSchema.FindColumn(column.Name)

			// Drop primary constraint.
			if curColumn.IsPrimaryKey() && !column.IsPrimaryKey() {
				err := db.DropTablePrimaryKey(
					schema.Name,
				)
				if err != nil {
					return err
				}
			}

			// Set primary constraint.
			if column.IsPrimaryKey() && !curColumn.IsPrimaryKey() {
				_, err := db.Exec(fmt.Sprintf(
					"alter table %s add primary key (%s)",
					schema.Name,
					column.Name,
				))
				if err != nil {
					return err
				}
			}

			// Type.
			columnBuf := column
			curColumnBuf := *curColumn

			columnBuf.Key.Type = NotKeyType
			curColumnBuf.Key.Type = NotKeyType

			columnSql, err := db.ColumnToSql(columnBuf)
			if err != nil {
				return err
			}

			curColumnSql, err := db.ColumnToSql(curColumnBuf)
			if err != nil {
				return err
			}

			if columnSql != curColumnSql {
				fmt.Printf("'%s'\n'%s'\n", columnSql, curColumnSql)
				_, err = db.Exec(fmt.Sprintf(
					"alter table %s modify column %s",
					schema.Name,
					columnSql,
				))

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

