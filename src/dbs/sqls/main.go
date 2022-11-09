package sqls

import (
	"fmt"
	"database/sql"
)

type DB struct {
	*sql.DB
	Driver string
	ConnConfig ConnConfig
}

type ConnConfig struct {
	Login, Password, Host, Name string
	Port int
}
type TableSchema struct {
	Name string
	Other string
}

func (c ConnConfig)String() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.Login,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
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

	for rows.Next() {
		s := TableSchema{}

		rows.Scan(
			&s.Name,
		)

		ret = append(ret, s)
	}

	return ret, nil
}

func Open(driver string, cfg ConnConfig) (*DB, error) {
	db, err := sql.Open(driver, cfg.String())
	if err != nil {
		return nil, err
	}

	return &DB{db, driver, cfg}, nil
}


