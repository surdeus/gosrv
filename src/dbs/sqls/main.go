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

func Open(driver string, cfg ConnConfig) (*DB, error) {
	db, err := sql.Open(driver, cfg.String())
	if err != nil {
		return nil, err
	}

	return &DB{db, driver, cfg}, nil
}


