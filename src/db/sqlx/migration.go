package sqlx

import (
	"fmt"
)

// Simple migration. Fuck writing much of shit in SQL directly.

func (db *DB)Migrate(sqlers []Sqler) error {
	var err error

	/*curSchemas, err := db.GetTableSchemas()
	if err != nil {
		return err
	}*/
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
			continue
		}

		db.CreateTableBySchema(schema)
	}

	// Then we rename existing and create not existing fields.
	for _, schema := range newSchemas {
	}

	return nil
}

