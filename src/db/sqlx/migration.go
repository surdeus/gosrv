package sqlx

import (
	"fmt"
)

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

	// First we should rename existing tables and create not existing ones.
	for _, schema := range newSchemas {
		// Rename.
		if schema.OldName != "" &&  db.TableExists(schema.OldName) {
			_, err = db.Query(fmt.Sprintf(
				"alter table %s rename %s ;",
				schema.OldName,
				schema.Name,
			))
			if err != nil {
				return err
			}

			// Fit changes to the current schema representation.
			curSchemas[curSchemas.FindSchema(schema.OldName)].Name = schema.Name

			continue
		}

		db.CreateTableBySchema(schema)
	}

	// Then we rename existing and create not existing fields.
	for _, schema := range newSchemas {
		idx := curSchemas.FindSchema(schema.Name)
		for _, field := range schema.Fields {
			if field.OldName != "" && db.FieldExists(schema.Name, field.OldName) {
				// Rename.
				_, err = db.Query(fmt.Sprintf(
					"alter table %s rename column %s to %s",
					schema.Name,
					field.OldName,
					field.Name,
				))
				if err != nil {
					return err
				}

				curFieldIdx := curSchemas[idx].FindField(field.OldName)
				fmt.Println(curFieldIdx)
				curSchemas[idx].Fields[curFieldIdx].Name = field.Name
			} else if !db.FieldExists(schema.Name, field.Name) {
				// Create.
				_, err = db.Query(fmt.Sprintf(
					"alter table %s add %s",
					schema.Name,
					db.FieldToSql(field),
				))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

